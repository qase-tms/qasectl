package xctest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

// failedTestActivities mirrors a failed XCUITest: every activity carries a failureSummaryIDs
// list, and only the failed one references a failure.
const failedTestActivities = `{"_values": [
  {
    "activityType": {"_value": "com.apple.dt.xctest.activity-type.userCreated"},
    "title": {"_value": "Tap login"},
    "failureSummaryIDs": {"_values": []},
    "attachments": {"_values": [
      {"filename": {"_value": "Debug description_0_A"}, "uniformTypeIdentifier": {"_value": "public.plain-text"}, "payloadRef": {"id": {"_value": "passed-top"}}}
    ]},
    "subactivities": {"_values": [
      {
        "activityType": {"_value": "com.apple.dt.xctest.activity-type.internal"},
        "title": {"_value": "Synthesize event"},
        "failureSummaryIDs": {"_values": []},
        "attachments": {"_values": [
          {"filename": {"_value": "SynthesizedEvent_1_B"}, "uniformTypeIdentifier": {"_value": "public.plain-text"}, "payloadRef": {"id": {"_value": "passed-nested"}}}
        ]}
      }
    ]}
  },
  {
    "activityType": {"_value": "com.apple.dt.xctest.activity-type.userCreated"},
    "title": {"_value": "Check title"},
    "failureSummaryIDs": {"_values": [{"_value": "f1"}]},
    "attachments": {"_values": [
      {"filename": {"_value": "App UI hierarchy for com.example_0_C"}, "uniformTypeIdentifier": {"_value": "public.plain-text"}, "payloadRef": {"id": {"_value": "failed-own"}}}
    ]}
  }
]}`

func newFailedTestParser(t *testing.T, mode string) (*Parser, ActivitySummaries) {
	t.Helper()
	p, err := NewParser("test.xcresult", "", mode)
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubXcrun(t, p, "echo text\n")

	var as ActivitySummaries
	if err := json.Unmarshal([]byte(failedTestActivities), &as); err != nil {
		t.Fatalf("failed to decode activities: %v", err)
	}

	var failure FailureSummary
	if err := json.Unmarshal([]byte(`{
		"uuid": {"_value": "f1"},
		"message": {"_value": "title mismatch"},
		"lineNumber": {"_value": "42"},
		"attachments": {"_values": [
			{"filename": {"_value": "Failure snapshot_0_D"}, "uniformTypeIdentifier": {"_value": "public.plain-text"}, "payloadRef": {"id": {"_value": "failure"}}}
		]}
	}`), &failure); err != nil {
		t.Fatalf("failed to decode failure: %v", err)
	}

	p.failures = map[string]FailureSummary{"f1": failure}
	p.processedFailures = nil
	p.stepPayloads = map[string]struct{}{}

	return p, as
}

func attachmentNames(attachments []models.Attachment) []string {
	names := make([]string, 0, len(attachments))
	for _, a := range attachments {
		names = append(names, a.Name)
	}
	return names
}

func equalNames(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestParser_getSteps_StepAttachmentModes(t *testing.T) {
	tests := []struct {
		mode         string
		passedTop    []string
		passedNested []string
		failed       []string
		resultLevel  []string
	}{
		{
			mode:         "",
			passedTop:    []string{},
			passedNested: []string{},
			failed:       []string{"Failure snapshot_0_D.txt", "App UI hierarchy for com.example_0_C.txt"},
			resultLevel:  []string{"Debug description_0_A.txt"},
		},
		{
			mode:         "failed",
			passedTop:    []string{},
			passedNested: []string{},
			failed:       []string{"Failure snapshot_0_D.txt", "App UI hierarchy for com.example_0_C.txt"},
			resultLevel:  []string{"Debug description_0_A.txt"},
		},
		{
			mode:         "all",
			passedTop:    []string{"Debug description_0_A.txt"},
			passedNested: []string{"SynthesizedEvent_1_B.txt"},
			failed:       []string{"Failure snapshot_0_D.txt", "App UI hierarchy for com.example_0_C.txt"},
			resultLevel:  []string{},
		},
		{
			mode:         "none",
			passedTop:    []string{},
			passedNested: []string{},
			failed:       []string{"Failure snapshot_0_D.txt"},
			resultLevel:  []string{"Debug description_0_A.txt", "App UI hierarchy for com.example_0_C.txt"},
		},
	}

	for _, tt := range tests {
		t.Run("mode="+tt.mode, func(t *testing.T) {
			p, as := newFailedTestParser(t, tt.mode)

			steps, failed, _ := p.getSteps(as, 0)
			if !failed {
				t.Error("getSteps() reported no failed step")
			}
			if len(steps) != 2 {
				t.Fatalf("got %d steps, want 2", len(steps))
			}

			if steps[0].Execution.Status != "passed" || steps[1].Execution.Status != "failed" {
				t.Errorf("statuses = %s, %s; want passed, failed", steps[0].Execution.Status, steps[1].Execution.Status)
			}
			if got := attachmentNames(steps[0].Execution.Attachments); !equalNames(got, tt.passedTop) {
				t.Errorf("passed step attachments = %q, want %q", got, tt.passedTop)
			}
			if got := attachmentNames(steps[0].Steps[0].Execution.Attachments); !equalNames(got, tt.passedNested) {
				t.Errorf("passed nested step attachments = %q, want %q", got, tt.passedNested)
			}
			if got := attachmentNames(steps[1].Execution.Attachments); !equalNames(got, tt.failed) {
				t.Errorf("failed step attachments = %q, want %q", got, tt.failed)
			}
			if got := attachmentNames(p.getTestLevelAttachments(as)); !equalNames(got, tt.resultLevel) {
				t.Errorf("result-level attachments = %q, want %q", got, tt.resultLevel)
			}
		})
	}
}

func TestNewParser_StepAttachmentModes(t *testing.T) {
	for _, mode := range []string{"", "failed", "FAILED", "all", "none"} {
		if _, err := NewParser("test.xcresult", "", mode); err != nil {
			t.Errorf("NewParser(mode=%q) error = %v", mode, err)
		}
	}
	if _, err := NewParser("test.xcresult", "", "passed"); err == nil {
		t.Error("NewParser(mode=passed) error = nil, want an error")
	}
}

// heicPayload is the start of a real HEIC file: the ftyp box with the heic brand.
const heicPayload = `printf '\000\000\000\030ftypheic\000\000\000\000'`

func stubSips(t *testing.T, p *Parser, body string) {
	t.Helper()
	p.sips = filepath.Join(t.TempDir(), "sips")
	if err := os.WriteFile(p.sips, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatalf("Failed to write sips stub: %v", err)
	}
}

func heicAttachment() Attachment {
	var a Attachment
	a.Filename.Value = "Screenshot_0_E.heic"
	a.UniformTypeIdentifier.Value = "public.heic"
	a.PayloadRef.ID.Value = "screenshot"
	return a
}

func TestParser_getAttachment_ConvertsHEICToJPEG(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubXcrun(t, p, heicPayload+"\n")
	// The last argument after --out is the destination file.
	stubSips(t, p, `for last; do :; done
printf '\377\330\377\340jpeg' > "$last"
`)

	got, err := p.getAttachment(heicAttachment())
	if err != nil {
		t.Fatalf("getAttachment() error = %v", err)
	}

	if got.Name != "Screenshot_0_E.jpeg" {
		t.Errorf("name = %q, want Screenshot_0_E.jpeg", got.Name)
	}
	if p.detectFileExtension(*got.Content) != ".jpg" {
		t.Errorf("content = %q, want JPEG data", *got.Content)
	}
}

func TestParser_getAttachment_KeepsHEICWhenConversionFails(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubXcrun(t, p, heicPayload+"\n")
	stubSips(t, p, "echo 'cannot convert' >&2\nexit 1\n")

	got, err := p.getAttachment(heicAttachment())
	if err != nil {
		t.Fatalf("getAttachment() error = %v", err)
	}

	if got.Name != "Screenshot_0_E.heic" {
		t.Errorf("name = %q, want Screenshot_0_E.heic", got.Name)
	}
	if p.detectFileExtension(*got.Content) != heicExt {
		t.Errorf("content = %q, want the original HEIC data", *got.Content)
	}
}

func TestParser_getAttachment_ConvertsHEIFBrandsByType(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	// mif1 is a generic HEIF brand that the content sniffing does not treat as HEIC.
	stubXcrun(t, p, `printf '\000\000\000\030ftypmif1\000\000\000\000'`+"\n")
	stubSips(t, p, `for last; do :; done
printf '\377\330\377\340jpeg' > "$last"
`)

	got, err := p.getAttachment(heicAttachment())
	if err != nil {
		t.Fatalf("getAttachment() error = %v", err)
	}

	if got.Name != "Screenshot_0_E.jpeg" {
		t.Errorf("name = %q, want Screenshot_0_E.jpeg", got.Name)
	}
}

func TestParser_getTestLevelAttachments_SkipsQaseConfig(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubXcrun(t, p, `echo '{"caseId":234}'`+"\n")
	p.stepPayloads = map[string]struct{}{}

	var as ActivitySummaries
	if err := json.Unmarshal([]byte(`{"_values": [{
		"activityType": {"_value": "com.apple.dt.xctest.activity-type.attachmentContainer"},
		"title": {"_value": "Added attachment named 'Qase config'"},
		"attachments": {"_values": [
			{"name": {"_value": "Qase config"}, "filename": {"_value": "Qase config_0_H"}, "uniformTypeIdentifier": {"_value": "public.data"}, "payloadRef": {"id": {"_value": "config"}}}
		]}
	}]}`), &as); err != nil {
		t.Fatalf("failed to decode activities: %v", err)
	}

	if got := p.getTestLevelAttachments(as); len(got) != 0 {
		t.Errorf("result-level attachments = %q, want none", attachmentNames(got))
	}
}

func TestParser_withExtension(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}

	tests := []struct {
		name     string
		filename string
		uti      string
		content  string
		want     string
	}{
		{name: "dot inside the name is not an extension", filename: "App UI hierarchy for com.example_0_C", uti: "public.plain-text", content: "Application", want: "App UI hierarchy for com.example_0_C.txt"},
		{name: "trailing dot", filename: "App UI hierarchy for com.apple_0_F.", uti: "public.plain-text", content: "Application", want: "App UI hierarchy for com.apple_0_F.txt"},
		{name: "matching extension is kept", filename: "Debug description_0_A.txt", uti: "public.plain-text", content: "Query", want: "Debug description_0_A.txt"},
		{name: "jpg counts as jpeg", filename: "photo.jpg", uti: "public.jpeg", content: "\xff\xd8\xff", want: "photo.jpg"},
		{name: "unknown type is sniffed", filename: "Snapshot_1_G", uti: "com.apple.dt.xctest.element-snapshot", content: "bplist00", want: "Snapshot_1_G.plist"},
		{name: "unknown type keeps a known extension", filename: "log.json", uti: "", content: "{}", want: "log.json"},
		{name: "unknown type and content stay as is", filename: "blob", uti: "", content: "\x00\x01\x02\x03", want: "blob"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.withExtension(tt.filename, tt.uti, []byte(tt.content)); got != tt.want {
				t.Errorf("withExtension(%q, %q) = %q, want %q", tt.filename, tt.uti, got, tt.want)
			}
		})
	}
}
