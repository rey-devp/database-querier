package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"database-querier-agent/pkg/llm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type QueryPlan struct {
	Collection string
	Operation  string
	Filter     bson.M
	Projection bson.M
	Pipeline   []bson.M
	Sort       bson.M
	Limit      int64
}

type Parser interface {
	Parse(ctx context.Context, userRequest string, availableCollections []string, sampleDocs map[string][]bson.M) (*QueryPlan, error)
}

type RuleBasedParser struct{}

func NewRuleBasedParser() *RuleBasedParser {
	return &RuleBasedParser{}
}

func (p *RuleBasedParser) Parse(ctx context.Context, userRequest string, availableCollections []string, sampleDocs map[string][]bson.M) (*QueryPlan, error) {
	req := strings.ToLower(strings.TrimSpace(userRequest))
	plan := &QueryPlan{
		Filter:     bson.M{},
		Projection: bson.M{},
	}

	// 1. Detect collection
	plan.Collection = detectCollection(req, availableCollections)

	// 2. Detect Operation
	plan.Operation = detectOperation(req)

	// 3. Extract Filters
	plan.Filter = extractFilters(req)

	// 4. Extract Projection
	plan.Projection = extractProjection(req)

	// 5. Aggregate pipeline (if needed)
	if plan.Operation == "aggregate" {
		plan.Pipeline = buildAggregatePipeline(req, plan.Filter)
		// Usually in aggregate, we don't pass filter separately, it goes into $match
	}

	return plan, nil
}

func detectCollection(req string, collections []string) string {
	// Simple substring match for now. "mahasiswa" -> "students" alias could be added.
	// Hardcoding an alias for this example since it's commonly asked in Indonesian
	if strings.Contains(req, "mahasiswa") {
		return "students"
	}
	
	for _, coll := range collections {
		if strings.Contains(req, coll) {
			return coll
		}
	}
	
	if len(collections) > 0 {
		return collections[0] // fallback
	}
	return "unknown"
}

func detectOperation(req string) string {
	for _, kw := range CountKeywords {
		if strings.Contains(req, kw) && !strings.Contains(req, "rata-rata") {
			return "countDocuments"
		}
	}
	
	for kw := range AggregationKeywords {
		if strings.Contains(req, kw) {
			return "aggregate"
		}
	}

	return "find"
}

func extractFilters(req string) bson.M {
	filter := bson.M{}
	
	// Example rule: "semester 6" -> {"semester": 6}
	reSemester := regexp.MustCompile(`semester\s+(\d+)`)
	if match := reSemester.FindStringSubmatch(req); match != nil {
		val, _ := strconv.Atoi(match[1])
		filter["semester"] = val
	}

	// Example rule: "gpa lebih dari 3.5"
	reGPA := regexp.MustCompile(`gpa\s+(lebih dari|kurang dari|sama dengan)\s+([0-9.]+)`)
	if match := reGPA.FindStringSubmatch(req); match != nil {
		op := OperatorMappings[match[1]]
		val, _ := strconv.ParseFloat(match[2], 64)
		
		if op == "$eq" {
			filter["gpa"] = val
		} else {
			filter["gpa"] = bson.M{op: val}
		}
	}

	return filter
}

func extractProjection(req string) bson.M {
	proj := bson.M{}
	if strings.Contains(req, "tampilkan nama dan gpa") {
		proj["name"] = 1
		proj["gpa"] = 1
	}
	return proj
}

func buildAggregatePipeline(req string, filter bson.M) []bson.M {
	var pipeline []bson.M

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	groupStage := bson.M{"_id": nil}
	
	if strings.Contains(req, "rata-rata gpa") {
		groupStage["avg"] = bson.M{"$avg": "$gpa"}
	}
	
	pipeline = append(pipeline, bson.M{"$group": groupStage})
	
	return pipeline
}

// --- LLM PARSER ---

type LLMParser struct {
	client llm.LLMClient
}

func NewLLMParser(client llm.LLMClient) *LLMParser {
	return &LLMParser{client: client}
}

func (p *LLMParser) Parse(ctx context.Context, userRequest string, availableCollections []string, sampleDocs map[string][]bson.M) (*QueryPlan, error) {
	prompt := llm.BuildPrompt(userRequest, availableCollections, sampleDocs)
	
	resp, err := p.client.GenerateQuery(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	plan := &QueryPlan{
		Collection: resp.Collection,
		Operation:  resp.Operation,
	}

	if resp.Filter != nil {
		plan.Filter = resp.Filter
	} else {
		plan.Filter = bson.M{}
	}

	if resp.Projection != nil {
		plan.Projection = resp.Projection
	} else {
		plan.Projection = bson.M{}
	}

	if resp.Sort != nil {
		plan.Sort = resp.Sort
	}

	if resp.Pipeline != nil {
		var pipeline []bson.M
		for _, stage := range resp.Pipeline {
			pipeline = append(pipeline, stage)
		}
		plan.Pipeline = pipeline
	}

	plan.Limit = resp.Limit

	return plan, nil
}
