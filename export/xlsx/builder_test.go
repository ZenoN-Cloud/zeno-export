package xlsx

import (
	"encoding/json"
	"os"
	"testing"
)

func TestBuildFromNormalized(t *testing.T) {
	// Read test data
	data, err := os.ReadFile("../../testdata/sample_normalized.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	var normalized map[string]interface{}
	if err := json.Unmarshal(data, &normalized); err != nil {
		t.Fatalf("failed to parse test data: %v", err)
	}

	// Generate Excel
	excelBytes, err := BuildFromNormalized(normalized)
	if err != nil {
		t.Fatalf("failed to build Excel: %v", err)
	}

	// Verify we got some data
	if len(excelBytes) == 0 {
		t.Error("expected Excel data, got empty bytes")
	}

	// Verify it's a valid ZIP (Excel is ZIP-based)
	if len(excelBytes) < 4 || string(excelBytes[:2]) != "PK" {
		t.Error("generated file doesn't appear to be a valid ZIP/Excel file")
	}

	t.Logf("Generated Excel file: %d bytes", len(excelBytes))
}

func TestBuildFromNormalized_EmptyData(t *testing.T) {
	normalized := map[string]interface{}{
		"vendor":   "test_bank",
		"rows":     []interface{}{},
		"warnings": []interface{}{},
	}

	excelBytes, err := BuildFromNormalized(normalized)
	if err != nil {
		t.Fatalf("failed to build Excel: %v", err)
	}

	if len(excelBytes) == 0 {
		t.Error("expected Excel data even for empty transactions")
	}
}