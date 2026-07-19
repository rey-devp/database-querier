package mongodb

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ReadOnlyOperations defines the allowed operation types
var ReadOnlyOperations = map[string]bool{
	"find":           true,
	"aggregate":      true,
	"countDocuments": true,
}

func ValidateOperation(operation string) error {
	if operation == "" {
		return fmt.Errorf("operation is empty")
	}

	if !ReadOnlyOperations[operation] {
		return fmt.Errorf("operation not allowed: '%s'. only read-only operations (find, aggregate, countDocuments) are permitted", operation)
	}

	return nil
}

func ValidatePipeline(pipeline []map[string]interface{}) error {
	for _, stage := range pipeline {
		for key := range stage {
			if strings.EqualFold(key, "$out") || strings.EqualFold(key, "$merge") {
				return fmt.Errorf("aggregation stage not allowed: '%s'. write operations are prohibited", key)
			}
		}
	}
	return nil
}

// ConvertToMapSlice converts a generic interface to []map[string]interface{} for validation
func ConvertToMapSlice(input interface{}) ([]map[string]interface{}, error) {
	if input == nil {
		return nil, nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(b, &result)
	return result, err
}
