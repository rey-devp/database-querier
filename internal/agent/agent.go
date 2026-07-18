package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"database-querier-agent/internal/logger"
	"database-querier-agent/internal/memory"
	"database-querier-agent/internal/mongodb"
	"database-querier-agent/internal/parser"
)

type Agent struct {
	parser *parser.RuleBasedParser
	mongo  *mongodb.Client
	store  *memory.Store
}

func NewAgent(mongo *mongodb.Client, store *memory.Store) *Agent {
	return &Agent{
		parser: parser.NewRuleBasedParser(),
		mongo:  mongo,
		store:  store,
	}
}

func (a *Agent) ProcessTask(ctx context.Context, taskID string) (*memory.AgentResponse, error) {
	startTime := time.Now()
	
	// 1. Get task
	task, err := a.store.GetTask(taskID)
	if err != nil {
		return nil, a.handleError(taskID, "Task not found")
	}
	
	logger.Info("AGENT", "Processing task", "task_id", taskID, "request", task.UserRequest)

	// 2. List collections
	collections, err := a.mongo.ListCollections(ctx)
	if err != nil {
		return nil, a.handleError(taskID, "Failed to list collections: "+err.Error())
	}
	logger.Info("AGENT", "Collections discovered", "count", len(collections))

	// 3. Parse user request
	plan, err := a.parser.Parse(task.UserRequest, collections)
	if err != nil {
		return nil, a.handleError(taskID, "Failed to parse request: "+err.Error())
	}
	logger.Info("PARSER", "Query plan generated", "collection", plan.Collection, "operation", plan.Operation)

	// 4. Validate operation (read-only)
	if err := mongodb.ValidateOperation(plan.Operation); err != nil {
		logger.Warn("VALIDATOR", "Operation REJECTED", "task_id", taskID, "attempted", plan.Operation)
		return nil, a.handleError(taskID, "Operation not allowed: "+err.Error())
	}
	logger.Info("VALIDATOR", "Read-only check passed", "operation", plan.Operation)

	// 5. Execute query
	coll := a.mongo.GetDatabase().Collection(plan.Collection)
	var documents interface{}
	var total int

	dbStartTime := time.Now()
	logger.Info("MONGODB", "Executing query...")

	switch plan.Operation {
	case "find":
		docs, count, err := mongodb.ExecuteFind(ctx, coll, plan.Filter, plan.Projection, plan.Sort, plan.Limit)
		if err != nil {
			return nil, a.handleError(taskID, "Find failed: "+err.Error())
		}
		documents = docs
		total = count
	case "aggregate":
		// Quick validation for pipeline
		pipelineMaps, _ := mongodb.ConvertToMapSlice(plan.Pipeline)
		if err := mongodb.ValidatePipeline(pipelineMaps); err != nil {
			logger.Warn("VALIDATOR", "Pipeline validation failed", "task_id", taskID, "error", err.Error())
			return nil, a.handleError(taskID, "Invalid pipeline: "+err.Error())
		}

		docs, count, err := mongodb.ExecuteAggregate(ctx, coll, plan.Pipeline)
		if err != nil {
			return nil, a.handleError(taskID, "Aggregate failed: "+err.Error())
		}
		documents = docs
		total = count
	case "countDocuments":
		docs, count, err := mongodb.ExecuteCountDocuments(ctx, coll, plan.Filter)
		if err != nil {
			return nil, a.handleError(taskID, "Count failed: "+err.Error())
		}
		documents = docs
		total = count
	default:
		logger.Warn("VALIDATOR", "Unsupported operation", "task_id", taskID, "attempted", plan.Operation)
		return nil, a.handleError(taskID, fmt.Sprintf("Unsupported operation: %s", plan.Operation))
	}
	
	dbDuration := time.Since(dbStartTime)
	logger.Info("MONGODB", "Query executed", "collection", plan.Collection, "total", total, "duration", dbDuration.String())

	// 6. Build response
	res := &memory.AgentResponse{
		AgentName: "database_querier",
		Result:    formatResult(plan.Operation, documents, total),
	}

	// 7. Save result
	a.store.SaveResult(taskID, res)
	
	totalDuration := time.Since(startTime)
	logger.Info("RESPONSE", "Task completed", "task_id", taskID, "status", "success", "duration", totalDuration.String())

	return res, nil
}

func (a *Agent) handleError(taskID, message string) error {
	errRes := &memory.ErrorResponse{
		AgentName: "database_querier",
		Result:    fmt.Sprintf("Gagal memproses permintaan: %s", message),
	}
	a.store.SaveResult(taskID, errRes)
	
	logger.Error("RESPONSE", "Task failed", "task_id", taskID, "message", message)
	
	return fmt.Errorf("%s", message)
}

func formatResult(operation string, documents interface{}, total int) string {
	if operation == "countDocuments" {
		return fmt.Sprintf("Total data yang ditemukan: %d", total)
	}

	if total == 0 {
		return "Tidak ditemukan data yang sesuai."
	}

	var sb strings.Builder
	if operation == "aggregate" {
		sb.WriteString("Hasil agregasi:\n")
	} else {
		sb.WriteString(fmt.Sprintf("Ditemukan %d data:\n", total))
	}

	docs, ok := documents.([]bson.M)
	if ok {
		for i, doc := range docs {
			var fields []string
			// Convert bson.M to something more readable if possible, or just JSON
			b, err := json.Marshal(doc)
			if err == nil {
				fields = append(fields, string(b))
			} else {
				for k, v := range doc {
					fields = append(fields, fmt.Sprintf("%s: %v", k, v))
				}
			}
			
			// If JSON was successful, it's just one element in `fields`
			if len(fields) == 1 && strings.HasPrefix(fields[0], "{") {
				sb.WriteString(fmt.Sprintf("- %s", fields[0]))
			} else {
				sb.WriteString(fmt.Sprintf("- %s", strings.Join(fields, ", ")))
			}
			
			if i < len(docs)-1 {
				sb.WriteString("\n")
			}
		}
	} else {
		// Fallback to JSON
		b, err := json.MarshalIndent(documents, "", "  ")
		if err == nil {
			sb.Write(b)
		} else {
			sb.WriteString(fmt.Sprintf("%v", documents))
		}
	}

	return sb.String()
}
