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
