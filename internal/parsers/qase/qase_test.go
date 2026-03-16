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

func TestParser_parseFile_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		wantErr bool
	}{
		{
			name:    "empty file",
			content: []byte{},
			wantErr: true,
		},
		{
			name:    "malformed JSON",
			content: []byte("{broken json"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.json")
			if err := os.WriteFile(tmpFile, tt.content, 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}
			parser := NewParser(tmpFile)
			results, err := parser.Parse()
			if tt.wantErr {
				// Qase parser logs errors and continues for dir-mode,
				// but for single-file mode errors propagate through parseFile
				if err != nil {
					return
				}
				if len(results) > 0 {
					t.Errorf("Parse() returned %d results for invalid input, want 0", len(results))
				}
				return
			}
			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}
		})
	}
}

func TestParser_Parse_FilenameWithSpaces(t *testing.T) {
	tmpDir := t.TempDir()
	spacedDir := filepath.Join(tmpDir, "dir with spaces")
	if err := os.MkdirAll(spacedDir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	jsonFile := filepath.Join(spacedDir, "test file.json")
	content := []byte(`{"title":"T","execution":{"status":"passed","duration":100}}`)
	if err := os.WriteFile(jsonFile, content, 0644); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	parser := NewParser(spacedDir)
	results, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() with spaces in path: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
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
