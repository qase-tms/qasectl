package maestro

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

func block(n int) []any {
	out := make([]any, n)
	for i := range out {
		out[i] = map[string]any{"tapOnElement": map[string]any{}}
	}
	return out
}

func actions(steps []models.Step) []string {
	out := make([]string, len(steps))
	for i, s := range steps {
		out[i] = s.Data.Action
	}
	return out
}

func TestBuildSteps_NestedSubflow(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "runFlowCommand", Status: "FAILED", Ts: 100, Dur: 50, Err: "Assertion is false: X"},
		{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 110, Dur: 20, Err: "Assertion is false: X"},
	})
	steps := buildSteps("", cmds, false)
	if len(steps) != 1 || len(steps[0].Steps) != 1 {
		t.Fatalf("tree: %+v", steps)
	}
	top := steps[0]
	if *top.Execution.StartTime != 100 || *top.Execution.EndTime != 150 || *top.Execution.Duration != 50 {
		t.Errorf("times: %v %v %v", *top.Execution.StartTime, *top.Execution.EndTime, *top.Execution.Duration)
	}
	child := top.Steps[0]
	if child.Execution.Status != "failed" || child.Execution.Comment != "Assertion is false: X" {
		t.Errorf("child: %+v", child.Execution)
	}
}

func TestBuildSteps_HidesServiceCommands(t *testing.T) {
	spec := append(header("H", "f", nil, nil, "U", "1"), cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 1100})
	cmds := decodeCommands(t, spec)
	if got := buildSteps("", cmds, false); len(got) != 1 {
		t.Errorf("hidden: %v", actions(got))
	}
	if got := buildSteps("", cmds, true); len(got) != 3 {
		t.Errorf("all: %v", actions(got))
	}
}

func TestBuildSteps_StatusMapping(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(prev)

	cases := []struct{ in, status, comment string }{
		{"COMPLETED", "passed", ""},
		{"FAILED", "failed", "boom"},
		{"SKIPPED", "skipped", ""},
		{"WARNED", "passed", "Optional command failed"},
		{"RUNNING", "failed", "Command did not finish"},
		{"PENDING", "failed", "Command did not finish"},
		{"WEIRD", "skipped", ""},
	}
	var spec []cmdSpec
	for _, c := range cases {
		s := cmdSpec{Type: "tapOnElement", Status: c.in, Ts: 1}
		if c.in == "FAILED" {
			s.Err = "boom"
		}
		spec = append(spec, s)
	}
	steps := buildSteps("", decodeCommands(t, spec), false)
	if len(steps) != len(cases) {
		t.Fatalf("got %d steps", len(steps))
	}
	for i, c := range cases {
		e := steps[i].Execution
		if e.Status != c.status || e.Comment != c.comment {
			t.Errorf("%s: got %q/%q, want %q/%q", c.in, e.Status, e.Comment, c.status, c.comment)
		}
	}
	if !strings.Contains(buf.String(), "WEIRD") {
		t.Errorf("no warning for unknown status: %q", buf.String())
	}
}

func TestBuildSteps_RetryRecovered(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "retryCommand", Args: map[string]any{"maxRetries": "3", "commands": block(3)}, Status: "COMPLETED", Ts: 100, Dur: 100},
		{Type: "evalScriptCommand", Status: "COMPLETED", Depth: 1, Ts: 101, Dur: 5},
		{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 110, Dur: 10, Err: "Assertion is false"},
		{Type: "evalScriptCommand", Status: "COMPLETED", Depth: 1, Ts: 130, Dur: 5},
		{Type: "assertConditionCommand", Status: "COMPLETED", Depth: 1, Ts: 140, Dur: 5},
		{Type: "assertConditionCommand", Status: "COMPLETED", Depth: 1, Ts: 150, Dur: 40},
	})
	steps := buildSteps("", cmds, false)
	if len(steps) != 1 {
		t.Fatalf("got %d top steps", len(steps))
	}
	att := steps[0].Steps
	if got := actions(att); len(got) != 2 || got[0] != "Attempt 1" || got[1] != "Attempt 2" {
		t.Fatalf("attempts: %v", got)
	}
	if att[0].Execution.Status != "failed" || att[1].Execution.Status != "passed" {
		t.Errorf("statuses: %q %q", att[0].Execution.Status, att[1].Execution.Status)
	}
	if len(att[0].Steps) != 2 || len(att[1].Steps) != 3 {
		t.Errorf("children: %d %d", len(att[0].Steps), len(att[1].Steps))
	}
	if *att[1].Execution.StartTime != 130 || *att[1].Execution.EndTime != 190 || *att[1].Execution.Duration != 60 {
		t.Errorf("attempt 2 times: %v %v %v", *att[1].Execution.StartTime, *att[1].Execution.EndTime, *att[1].Execution.Duration)
	}
	if !hadFailedAttempt(cmds) {
		t.Error("hadFailedAttempt = false")
	}
}

func TestBuildSteps_RetryExhausted(t *testing.T) {
	spec := []cmdSpec{{Type: "retryCommand", Args: map[string]any{"maxRetries": "3", "commands": block(1)}, Status: "FAILED", Ts: 100, Dur: 100, Err: "e"}}
	for i := 0; i < 4; i++ {
		spec = append(spec, cmdSpec{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: int64(101 + i*10), Dur: 5, Err: "e"})
	}
	att := buildSteps("", decodeCommands(t, spec), false)[0].Steps
	if len(att) != 4 {
		t.Fatalf("attempts: %v", actions(att))
	}
	for i, a := range att {
		if a.Execution.Status != "failed" || len(a.Steps) != 1 {
			t.Errorf("attempt %d: %q, %d children", i+1, a.Execution.Status, len(a.Steps))
		}
	}
}

func TestBuildSteps_RetryWithNestedRunFlow(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "retryCommand", Args: map[string]any{"maxRetries": "1", "commands": block(2)}, Status: "COMPLETED", Ts: 100, Dur: 100},
		{Type: "runFlowCommand", Status: "COMPLETED", Depth: 1, Ts: 101, Dur: 50},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 2, Ts: 102, Dur: 5},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 2, Ts: 110, Dur: 5},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 2, Ts: 120, Dur: 5},
		{Type: "assertConditionCommand", Status: "COMPLETED", Depth: 1, Ts: 160, Dur: 5},
	})
	att := buildSteps("", cmds, false)[0].Steps
	if len(att) != 1 || len(att[0].Steps) != 2 {
		t.Fatalf("attempts: %+v", att)
	}
	if len(att[0].Steps[0].Steps) != 3 {
		t.Errorf("runFlow children: %d", len(att[0].Steps[0].Steps))
	}
	if hadFailedAttempt(cmds) {
		t.Error("hadFailedAttempt = true")
	}
}

func TestBuildSteps_Repeat(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "repeatCommand", Args: map[string]any{"times": "2", "commands": block(1)}, Status: "COMPLETED", Ts: 100, Dur: 100},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 1, Ts: 101, Dur: 5},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 1, Ts: 110, Dur: 5},
	})
	att := buildSteps("", cmds, false)[0].Steps
	if got := actions(att); len(got) != 2 || got[0] != "Iteration 1" || got[1] != "Iteration 2" {
		t.Fatalf("iterations: %v", got)
	}
}

func TestBuildSteps_UnreadableBlockKeepsChildrenFlat(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "retryCommand", Args: map[string]any{"maxRetries": "1"}, Status: "COMPLETED", Ts: 100, Dur: 100},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 1, Ts: 101, Dur: 5},
		{Type: "tapOnElement", Status: "COMPLETED", Depth: 1, Ts: 110, Dur: 5},
	})
	kids := buildSteps("", cmds, false)[0].Steps
	if got := actions(kids); len(got) != 2 || strings.HasPrefix(got[0], "Attempt") {
		t.Errorf("children: %v", got)
	}
}

func TestDescribe(t *testing.T) {
	cases := []struct {
		name string
		spec cmdSpec
		want string
	}{
		{"label wins", cmdSpec{Type: "tapOnElement", Args: map[string]any{"label": "My label", "selector": map[string]any{"textRegex": "Skip"}}}, "My label"},
		{"tap", cmdSpec{Type: "tapOnElement", Args: map[string]any{"selector": map[string]any{"textRegex": "Skip", "optional": false}}}, `Tap on "Skip"`},
		{"assert visible", cmdSpec{Type: "assertConditionCommand", Args: map[string]any{"condition": map[string]any{"visible": map[string]any{"textRegex": "Search Wikipedia"}}}}, `Assert visible "Search Wikipedia"`},
		{"assert script", cmdSpec{Type: "assertConditionCommand", Args: map[string]any{"condition": map[string]any{"scriptCondition": "${output.attempt >= 2}"}}}, "Assert true: output.attempt >= 2"},
		{"run flow", cmdSpec{Type: "runFlowCommand", Args: map[string]any{"sourceDescription": "../../subflows/x.yaml"}}, "Run flow ../../subflows/x.yaml"},
		{"input", cmdSpec{Type: "inputTextCommand", Args: map[string]any{"text": "hello"}}, `Input text "hello"`},
		{"unknown", cmdSpec{Type: "fooBarCommand"}, "fooBar"},
	}
	for _, c := range cases {
		cmds := decodeCommands(t, []cmdSpec{c.spec})
		if got := describe(cmds[0]); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestDeepestFailure_FollowsLastFailedChain(t *testing.T) {
	cmds := decodeCommands(t, []cmdSpec{
		{Type: "retryCommand", Args: map[string]any{"maxRetries": "1", "commands": block(1)}, Status: "FAILED", Ts: 1, Err: "outer"},
		{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 2, Err: "first"},
		{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 3, Err: "last"},
	})
	if f := deepestFailure(cmds); f == nil || f.Error.Message != "last" {
		t.Errorf("got %+v", f)
	}
	if f := deepestFailure(nil); f != nil {
		t.Errorf("empty: %+v", f)
	}
}

func TestParse_FailureMessageFromDeepest(t *testing.T) {
	root := t.TempDir()
	cmds := append(okFlow("F"),
		cmdSpec{Type: "runFlowCommand", Status: "FAILED", Ts: 1200, Dur: 50, Err: "outer", Debug: "outer debug"},
		cmdSpec{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 1210, Dur: 5, Err: "inner", Debug: "inner debug"},
	)
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	r := res[0]
	if r.Message == nil || *r.Message != "inner" {
		t.Errorf("message %v", r.Message)
	}
	if r.Execution.StackTrace == nil || *r.Execution.StackTrace != "inner debug" {
		t.Errorf("stacktrace %v", r.Execution.StackTrace)
	}
	if len(r.Steps) != 2 {
		t.Errorf("steps: %v", actions(r.Steps))
	}
}

func TestParse_StepsAllShowsServiceCommands(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "x", okFlow("S"))
	res, err := NewParser(root, "ALL").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res[0].Steps) != 3 {
		t.Errorf("steps: %v", actions(res[0].Steps))
	}
}

func TestParse_FlakyField(t *testing.T) {
	retry := func(last string) []cmdSpec {
		return []cmdSpec{
			{Type: "retryCommand", Args: map[string]any{"maxRetries": "1", "commands": block(1)}, Status: last, Ts: 1200, Dur: 50},
			{Type: "assertConditionCommand", Status: "FAILED", Depth: 1, Ts: 1201, Dur: 5, Err: "e"},
			{Type: "assertConditionCommand", Status: last, Depth: 1, Ts: 1210, Dur: 5},
		}
	}
	cases := []struct {
		name  string
		extra []cmdSpec
		flaky bool
	}{
		{"recovered", retry("COMPLETED"), true},
		{"exhausted", retry("FAILED"), false},
		{"no failed attempt", nil, false},
	}
	for _, c := range cases {
		root := t.TempDir()
		writeFlow(t, root, "x", append(okFlow(c.name), c.extra...))
		res, err := NewParser(root, "").Parse()
		if err != nil {
			t.Fatal(err)
		}
		v, ok := res[0].Fields["is_flaky"]
		if c.flaky && v != "true" || !c.flaky && ok {
			t.Errorf("%s: is_flaky=%q present=%v", c.name, v, ok)
		}
	}
}

func TestParse_ExhaustedRetryInOptionalFlowNotFlaky(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "x", append(okFlow("opt"),
		cmdSpec{Type: "runFlowCommand", Status: "WARNED", Ts: 1200, Dur: 50},
		cmdSpec{Type: "retryCommand", Args: map[string]any{"maxRetries": "1", "commands": block(1)}, Status: "FAILED", Depth: 1, Ts: 1201, Dur: 40},
		cmdSpec{Type: "assertConditionCommand", Status: "FAILED", Depth: 2, Ts: 1202, Dur: 5, Err: "e"},
		cmdSpec{Type: "assertConditionCommand", Status: "FAILED", Depth: 2, Ts: 1210, Dur: 5, Err: "e"},
	))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Execution.Status != "passed" {
		t.Errorf("status: %q", res[0].Execution.Status)
	}
	if v, ok := res[0].Fields["is_flaky"]; ok {
		t.Errorf("is_flaky=%q must be absent", v)
	}
}
