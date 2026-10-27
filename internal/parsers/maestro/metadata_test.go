package maestro

import (
	"reflect"
	"testing"
)

func TestParseQaseIDs(t *testing.T) {
	tests := []struct {
		in   string
		want []int64
	}{
		{"123", []int64{123}},
		{"123, 456", []int64{123, 456}},
		{" PRJ-123 ,PRJ-456 ", []int64{123, 456}},
		{"123, abc", []int64{123}},
		{"", nil},
	}
	for _, tt := range tests {
		if got := parseQaseIDs(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseQaseIDs(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParseSuites(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"Search.Smoke", []string{"Search", "Smoke"}},
		{"", nil},
		{"A..B", []string{"A", "B"}},
	}
	for _, tt := range tests {
		if got := parseSuites(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseSuites(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestParse_Metadata(t *testing.T) {
	root := t.TempDir()
	props := map[string]string{"qase.id": "12, 15", "qase.suite": "Search.Smoke", "priority": "high"}
	cmds := header("03 · Search for an article", "03-search-article", []string{"smoke", "search"}, props, "U", "1")
	writeFlow(t, root, "a", append(cmds, cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 1100}))

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	r := res[0]
	if r.TestOpsIDs == nil || !reflect.DeepEqual(*r.TestOpsIDs, []int64{12, 15}) {
		t.Errorf("TestOpsIDs = %v", r.TestOpsIDs)
	}
	if r.TestOpsID != nil {
		t.Errorf("TestOpsID = %v, want nil", *r.TestOpsID)
	}
	var titles []string
	for _, s := range r.Relations.Suite.Data {
		titles = append(titles, s.Title)
	}
	if !reflect.DeepEqual(titles, []string{"Search", "Smoke"}) {
		t.Errorf("suites = %v", titles)
	}
	want := map[string]string{"priority": "high", "tags": "smoke, search"}
	if !reflect.DeepEqual(r.Fields, want) {
		t.Errorf("Fields = %v, want %v", r.Fields, want)
	}
	wantSig := "12-15::search::smoke::03-search-article::03_·_search_for_an_article"
	if r.Signature == nil || *r.Signature != wantSig {
		t.Errorf("Signature = %v, want %s", r.Signature, wantSig)
	}
}

func TestParse_NoTagsNoField(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "a", append(header("A", "f", nil, nil, "U", "1"),
		cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 1100}))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	r := res[0]
	if _, ok := r.Fields["tags"]; ok {
		t.Errorf("unexpected tags field: %v", r.Fields)
	}
	if r.Fields == nil {
		t.Error("Fields must be non-nil")
	}
	if r.TestOpsIDs != nil || len(r.Relations.Suite.Data) != 0 {
		t.Errorf("unexpected ids/suites: %v %v", r.TestOpsIDs, r.Relations)
	}
}
