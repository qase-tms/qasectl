package xctest

import (
	"testing"
)

// TestParser_InstanceIsolation verifies that Parser struct fields
// are instance-scoped, not shared across instances.
func TestParser_InstanceIsolation(t *testing.T) {
	// Create two parser instances (use invalid paths — we won't call Parse)
	p1, err := NewParser("test1.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser 1: %v", err)
	}
	p2, err := NewParser("test2.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser 2: %v", err)
	}

	// Initialize fields on p1
	var id1 int64 = 42
	p1.caseId = &id1
	p1.failures = map[string]FailureSummary{"key1": {}}
	p1.processedFailures = []string{"f1"}

	// Initialize fields on p2
	var id2 int64 = 99
	p2.caseId = &id2
	p2.failures = map[string]FailureSummary{"key2": {}}
	p2.processedFailures = []string{"f2"}

	// Verify p1 fields are not affected by p2
	if *p1.caseId != 42 {
		t.Errorf("p1.caseId = %d, want 42", *p1.caseId)
	}
	if _, ok := p1.failures["key1"]; !ok {
		t.Error("p1.failures missing key1")
	}
	if _, ok := p1.failures["key2"]; ok {
		t.Error("p1.failures unexpectedly has key2 from p2")
	}
	if len(p1.processedFailures) != 1 || p1.processedFailures[0] != "f1" {
		t.Errorf("p1.processedFailures = %v, want [f1]", p1.processedFailures)
	}

	// Verify p2 fields are not affected by p1
	if *p2.caseId != 99 {
		t.Errorf("p2.caseId = %d, want 99", *p2.caseId)
	}
	if _, ok := p2.failures["key2"]; !ok {
		t.Error("p2.failures missing key2")
	}
	if _, ok := p2.failures["key1"]; ok {
		t.Error("p2.failures unexpectedly has key1 from p1")
	}
}

// TestParser_IsFailureProcessed verifies the method works on instance state
func TestParser_IsFailureProcessed(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	p.processedFailures = []string{"abc", "def"}

	if !p.isFailureProcessed("abc") {
		t.Error("isFailureProcessed(abc) = false, want true")
	}
	if !p.isFailureProcessed("def") {
		t.Error("isFailureProcessed(def) = false, want true")
	}
	if p.isFailureProcessed("ghi") {
		t.Error("isFailureProcessed(ghi) = true, want false")
	}
}

// makeHeicData creates a minimal byte slice with the given ftyp signature at bytes 4-12
func makeHeicData(ftyp string) []byte {
	data := make([]byte, 16)
	// bytes 0-3: arbitrary (file size in real HEIC)
	data[0], data[1], data[2], data[3] = 0x00, 0x00, 0x00, 0x18
	copy(data[4:12], ftyp)
	return data
}

func TestParser_detectFileExtension(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}

	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "PNG signature",
			data:     []byte("\x89PNG\r\n\x1a\n" + "extra"),
			expected: ".png",
		},
		{
			name:     "JPG signature",
			data:     []byte("\xff\xd8\xff\xe0"),
			expected: ".jpg",
		},
		// All 7 HEIC ftyp signatures
		{name: "HEIC ftypheic", data: makeHeicData("ftypheic"), expected: ".heic"},
		{name: "HEIC ftypheix", data: makeHeicData("ftypheix"), expected: ".heic"},
		{name: "HEIC ftypheis", data: makeHeicData("ftypheis"), expected: ".heic"},
		{name: "HEIC ftyphevc", data: makeHeicData("ftyphevc"), expected: ".heic"},
		{name: "HEIC ftyphevx", data: makeHeicData("ftyphevx"), expected: ".heic"},
		{name: "HEIC ftyphevs", data: makeHeicData("ftyphevs"), expected: ".heic"},
		{name: "HEIC ftyphevm", data: makeHeicData("ftyphevm"), expected: ".heic"},
		{
			name:     "PDF signature",
			data:     []byte("%PDF-1.4"),
			expected: ".pdf",
		},
		{
			name:     "JSON object",
			data:     []byte(`{"key": "value"}`),
			expected: ".json",
		},
		{
			name:     "unknown bytes",
			data:     []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D},
			expected: "",
		},
		{
			name:     "too short",
			data:     []byte{0x01, 0x02},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.detectFileExtension(tt.data)
			if got != tt.expected {
				t.Errorf("detectFileExtension() = %q, want %q", got, tt.expected)
			}
		})
	}
}
