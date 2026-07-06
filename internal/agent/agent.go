package agent

import (
	"context"
	"fmt"
	"time"

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
		MemoryPayload: memory.MemoryPayload{
			Result: memory.QueryResult{
				Collection: plan.Collection,
				Operation:  plan.Operation,
				Filter:     plan.Filter,
				Projection: plan.Projection,
				Documents:  documents,
				Total:      total,
			},
		},
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
		Status:    "failed",
		Message:   message,
	}
	a.store.SaveResult(taskID, errRes)
	
	logger.Error("RESPONSE", "Task failed", "task_id", taskID, "message", message)
	
	return fmt.Errorf("%s", message)
}
