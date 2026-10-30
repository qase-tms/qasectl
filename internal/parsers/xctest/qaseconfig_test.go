package xctest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

func int64Ptr(v int64) *int64 { return &v }

func TestQaseConfig_merge(t *testing.T) {
	tests := []struct {
		name    string
		configs []QaseConfig
		want    QaseConfig
	}{
		{
			name:    "deprecated caseId becomes an ID",
			configs: []QaseConfig{{CaseID: int64Ptr(7)}},
			want:    QaseConfig{IDs: []int64{7}},
		},
		{
			name:    "caseId and ids are combined without duplicates",
			configs: []QaseConfig{{CaseID: int64Ptr(7), IDs: []int64{7, 8}}, {IDs: []int64{8, 9}}},
			want:    QaseConfig{IDs: []int64{7, 8, 9}},
		},
		{
			name: "later title and suite win, empty ones are ignored",
			configs: []QaseConfig{
				{Title: "first", Suite: "A.B"},
				{Title: "second"},
				{Suite: "C"},
			},
			want: QaseConfig{Title: "second", Suite: "C"},
		},
		{
			name: "fields and params add up, later keys win",
			configs: []QaseConfig{
				{Fields: map[string]string{"severity": "minor", "layer": "e2e"}, Params: map[string]string{"user": "guest"}},
				{Fields: map[string]string{"severity": "critical"}, Params: map[string]string{"locale": "en"}},
			},
			want: QaseConfig{
				Fields: map[string]string{"severity": "critical", "layer": "e2e"},
				Params: map[string]string{"user": "guest", "locale": "en"},
			},
		},
		{
			name:    "ignore and muted stay set once any config sets them",
			configs: []QaseConfig{{Ignore: true, Muted: true}, {}},
			want:    QaseConfig{Ignore: true, Muted: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got QaseConfig
			for _, c := range tt.configs {
				got.merge(c)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestQaseConfig_apply(t *testing.T) {
	base := func() models.Result {
		return models.Result{
			Title:  "testLogin()",
			Fields: map[string]string{},
			Params: map[string]string{"Device": "iPhone 16", "Configuration": "Default"},
			Relations: models.Relation{Suite: models.Suite{Data: []models.SuiteData{
				{Title: ""}, {Title: "AppUITests"}, {Title: "LoginTests"},
			}}},
		}
	}

	t.Run("empty config keeps the defaults", func(t *testing.T) {
		got := base()
		QaseConfig{}.apply(&got)
		if !reflect.DeepEqual(got, base()) {
			t.Errorf("apply() = %+v, want %+v", got, base())
		}
	})

	t.Run("one ID sets only the case ID", func(t *testing.T) {
		got := base()
		QaseConfig{IDs: []int64{5}}.apply(&got)
		if got.TestOpsID == nil || *got.TestOpsID != 5 {
			t.Errorf("TestOpsID = %v, want 5", got.TestOpsID)
		}
		if got.TestOpsIDs != nil {
			t.Errorf("TestOpsIDs = %v, want nil", *got.TestOpsIDs)
		}
	})

	t.Run("full config overrides the result", func(t *testing.T) {
		got := base()
		QaseConfig{
			IDs:    []int64{5, 6},
			Title:  "Login with valid credentials",
			Suite:  " Auth . Login ",
			Fields: map[string]string{"severity": "critical"},
			Params: map[string]string{"user": "admin"},
			Muted:  true,
		}.apply(&got)

		if *got.TestOpsID != 5 || !reflect.DeepEqual(*got.TestOpsIDs, []int64{5, 6}) {
			t.Errorf("IDs = %v / %v, want 5 / [5 6]", *got.TestOpsID, *got.TestOpsIDs)
		}
		if got.Title != "Login with valid credentials" {
			t.Errorf("Title = %q", got.Title)
		}
		wantSuites := []models.SuiteData{{Title: "Auth"}, {Title: "Login"}}
		if !reflect.DeepEqual(got.Relations.Suite.Data, wantSuites) {
			t.Errorf("suites = %+v, want %+v", got.Relations.Suite.Data, wantSuites)
		}
		wantFields := map[string]string{"severity": "critical", "muted": "true"}
		if !reflect.DeepEqual(got.Fields, wantFields) {
			t.Errorf("Fields = %v, want %v", got.Fields, wantFields)
		}
		wantParams := map[string]string{"Device": "iPhone 16", "Configuration": "Default", "user": "admin"}
		if !reflect.DeepEqual(got.Params, wantParams) {
			t.Errorf("Params = %v, want %v", got.Params, wantParams)
		}
		if !got.Muted {
			t.Error("Muted = false, want true")
		}
	})
}

// configActivities builds the activities of a test that adds one Qase config attachment per
// payload ID, the way XCTAttachment records them.
func configActivities(t *testing.T, payloadIDs ...string) ActivitySummaries {
	t.Helper()
	values := make([]string, 0, len(payloadIDs)+1)
	values = append(values, `{"activityType": {"_value": "com.apple.dt.xctest.activity-type.userCreated"}, "title": {"_value": "Tap login"}}`)
	for _, id := range payloadIDs {
		values = append(values, fmt.Sprintf(`{
			"activityType": {"_value": "com.apple.dt.xctest.activity-type.attachmentContainer"},
			"title": {"_value": "Added attachment named 'Qase config'"},
			"attachments": {"_values": [
				{"name": {"_value": "Qase config"}, "filename": {"_value": "Qase config_0_%[1]s"}, "uniformTypeIdentifier": {"_value": "public.data"}, "payloadRef": {"id": {"_value": "%[1]s"}}}
			]}
		}`, id))
	}

	var as ActivitySummaries
	if err := json.Unmarshal([]byte(`{"_values": [`+strings.Join(values, ",")+`]}`), &as); err != nil {
		t.Fatalf("failed to decode activities: %v", err)
	}
	return as
}

// stubConfigs makes xcrun return the JSON mapped to the requested payload ID.
func stubConfigs(t *testing.T, p *Parser, payloads map[string]string) {
	t.Helper()
	var body strings.Builder
	body.WriteString("for last; do :; done\ncase \"$last\" in\n")
	for id, content := range payloads {
		fmt.Fprintf(&body, "%s) echo '%s' ;;\n", id, content)
	}
	body.WriteString("esac\n")
	stubXcrun(t, p, body.String())
}

func xcTest(as ActivitySummaries) XCTest {
	return XCTest{
		Name:      "testLogin()",
		Suites:    []string{"LoginTests"},
		Signature: "test://com.apple.xcode/App/AppUITests/LoginTests/testLogin",
		Action:    ActionTestSummary{ActivitySummaries: as},
	}
}

func TestParser_getResults_QaseConfig(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubConfigs(t, p, map[string]string{
		"cfg1": `{"ids":[11,12],"title":"Login works","fields":{"severity":"critical"}}`,
		"cfg2": `{"suite":"Auth.Login","params":{"user":"admin"},"fields":{"layer":"e2e"}}`,
	})

	results := p.getResults([]XCTest{xcTest(configActivities(t, "cfg1", "cfg2"))})
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	r := results[0]

	if r.TestOpsIDs == nil || !reflect.DeepEqual(*r.TestOpsIDs, []int64{11, 12}) {
		t.Errorf("TestOpsIDs = %v, want [11 12]", r.TestOpsIDs)
	}
	if r.Title != "Login works" {
		t.Errorf("Title = %q, want Login works", r.Title)
	}
	if want := []models.SuiteData{{Title: "Auth"}, {Title: "Login"}}; !reflect.DeepEqual(r.Relations.Suite.Data, want) {
		t.Errorf("suites = %+v, want %+v", r.Relations.Suite.Data, want)
	}
	if want := map[string]string{"severity": "critical", "layer": "e2e"}; !reflect.DeepEqual(r.Fields, want) {
		t.Errorf("Fields = %v, want %v", r.Fields, want)
	}
	if r.Params["user"] != "admin" {
		t.Errorf("Params = %v, want user=admin", r.Params)
	}
	if *r.Signature != "test://com.apple.xcode/App/AppUITests/LoginTests/testLogin" {
		t.Errorf("Signature = %q, the title must not change it", *r.Signature)
	}
	if len(r.Steps) != 1 || len(r.Attachments) != 0 {
		t.Errorf("steps = %d, attachments = %q; the configs must be neither", len(r.Steps), attachmentNames(r.Attachments))
	}
}

func TestParser_getResults_LegacyCaseID(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubConfigs(t, p, map[string]string{"cfg": `{"caseId":234}`})

	results := p.getResults([]XCTest{xcTest(configActivities(t, "cfg"))})
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if id := results[0].TestOpsID; id == nil || *id != 234 {
		t.Errorf("TestOpsID = %v, want 234", id)
	}
	if !p.caseIDWarned {
		t.Error("caseId must log a deprecation warning")
	}
}

func TestParser_getResults_ConfigIsPerTest(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubConfigs(t, p, map[string]string{
		"ignored": `{"ids":[1],"ignore":true}`,
		"kept":    `{"ids":[2]}`,
	})

	results := p.getResults([]XCTest{
		xcTest(configActivities(t, "ignored")),
		xcTest(configActivities(t, "kept")),
		xcTest(configActivities(t)),
	})
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2: the ignored test must be dropped", len(results))
	}
	if id := results[0].TestOpsID; id == nil || *id != 2 {
		t.Errorf("first result TestOpsID = %v, want 2", id)
	}
	if results[1].TestOpsID != nil {
		t.Errorf("a test without a config got TestOpsID %d from the previous one", *results[1].TestOpsID)
	}
}

func TestParser_getResults_BrokenConfigIsSkipped(t *testing.T) {
	p, err := NewParser("test.xcresult", "", "")
	if err != nil {
		t.Fatalf("NewParser() error = %v", err)
	}
	stubConfigs(t, p, map[string]string{
		"broken": `{"ids":"12"}`,
		"good":   `{"title":"Kept"}`,
	})

	results := p.getResults([]XCTest{xcTest(configActivities(t, "broken", "good"))})
	if len(results) != 1 || results[0].Title != "Kept" || results[0].TestOpsID != nil {
		t.Errorf("results = %+v, want one result titled Kept without IDs", results)
	}
}
