package qase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_WithBOM(t *testing.T) {
	content := []byte("\xef\xbb\xbf" + `{
		"title": "BOM Test",
		"execution": {
			"status": "passed",
			"duration": 100
		}
	}`)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "bom_test.json")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	parser := NewParser(filePath)
	results, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Title != "BOM Test" {
		t.Errorf("expected title 'BOM Test', got '%s'", results[0].Title)
	}
}

func TestParseFile_WithoutBOM(t *testing.T) {
	content := []byte(`{
		"title": "No BOM Test",
		"execution": {
			"status": "passed",
			"duration": 100
		}
	}`)

	dir := t.TempDir()
	filePath := filepath.Join(dir, "no_bom_test.json")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	parser := NewParser(filePath)
	results, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Title != "No BOM Test" {
		t.Errorf("expected title 'No BOM Test', got '%s'", results[0].Title)
	}
}
