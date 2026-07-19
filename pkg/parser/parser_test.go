package parser

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestParseOperation(t *testing.T) {
	parser := NewRuleBasedParser()
	collections := []string{"students"}

	tests := []struct {
		request  string
		expected string
	}{
		{"Tampilkan seluruh mahasiswa", "find"},
		{"Hitung jumlah mahasiswa", "countDocuments"},
		{"Berapa banyak data mahasiswa", "countDocuments"},
		{"Rata-rata gpa mahasiswa", "aggregate"},
	}

	for _, tc := range tests {
		plan, _ := parser.Parse(tc.request, collections)
		if plan.Operation != tc.expected {
			t.Errorf("Expected %s, got %s for request: %s", tc.expected, plan.Operation, tc.request)
		}
	}
}

func TestParseFilter(t *testing.T) {
	parser := NewRuleBasedParser()
	collections := []string{"students"}

	request := "Tampilkan mahasiswa semester 6 dengan gpa lebih dari 3.5"
	plan, _ := parser.Parse(request, collections)

	if plan.Filter["semester"] != 6 {
		t.Errorf("Expected semester to be 6, got %v", plan.Filter["semester"])
	}

	gpaFilter, ok := plan.Filter["gpa"].(bson.M)
	if !ok || gpaFilter["$gt"] != 3.5 {
		t.Errorf("Expected gpa filter to be {$gt: 3.5}, got %v", plan.Filter["gpa"])
	}
}
