package maestro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func okFlow(name string) []cmdSpec {
	cmds := header(name, "f", nil, nil, "UDID", "1")
	return append(cmds, cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 1100, Dur: 10})
}

func TestParse_OneResultPerFlow(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "a", okFlow("A"))
	bad := okFlow("B")
	bad = append(bad, cmdSpec{Type: "assertVisible", Status: "FAILED", Ts: 1200})
	writeFlow(t, root, "b", bad)

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("got %d results", len(res))
	}
	got := map[string]string{}
	for _, r := range res {
		got[r.Title] = r.Execution.Status
	}
	if got["A"] != "passed" || got["B"] != "failed" {
		t.Errorf("statuses: %v", got)
	}
}

func TestParse_TitleFromConfigNotFolder(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "03 · Search — поиск-shard-2", okFlow("03 · Search — поиск"))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Title != "03 · Search — поиск" {
		t.Errorf("title %q", res[0].Title)
	}
}

func TestParse_TitleFallsBackToFileName(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "x", header("", "03-search-article", nil, nil, "U", "1"))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Title != "03-search-article" {
		t.Errorf("title %q", res[0].Title)
	}
}

func TestParse_TimesFromCommands(t *testing.T) {
	root := t.TempDir()
	cmds := header("T", "f", nil, nil, "U", "1")
	cmds[0].Ts = 1000
	cmds = append(cmds, cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 5000, Dur: 250})
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := res[0].Execution
	if *e.StartTime != 1000 || *e.EndTime != 5250 || *e.Duration != 4250 {
		t.Errorf("start=%v end=%v dur=%v", *e.StartTime, *e.EndTime, *e.Duration)
	}
}

func TestParse_InterruptedFlowFails(t *testing.T) {
	root := t.TempDir()
	cmds := okFlow("I")
	cmds = append(cmds, cmdSpec{Type: "tapOnElement", Status: "RUNNING", Ts: 1300})
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Execution.Status != "failed" {
		t.Errorf("status %q", res[0].Execution.Status)
	}
	if res[0].Message == nil || *res[0].Message != "Flow did not finish" {
		t.Errorf("message %v", res[0].Message)
	}
}

func TestParse_SkipsEmptyOrBrokenFlow(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "empty", nil)
	broken := filepath.Join(root, "broken")
	if err := os.MkdirAll(broken, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(broken, "commands.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFlow(t, root, "nocfg", []cmdSpec{{Type: "tapOnElement", Status: "COMPLETED", Ts: 1}})
	writeFlow(t, root, "ok", okFlow("OK"))

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].Title != "OK" {
		t.Fatalf("results: %+v", res)
	}
}

func TestParse_PathIsSingleFlowDir(t *testing.T) {
	root := t.TempDir()
	dir := writeFlow(t, root, "only", okFlow("Only"))
	res, err := NewParser(dir, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("got %d", len(res))
	}
}

func TestParse_NoCommandsJSON(t *testing.T) {
	_, err := NewParser(t.TempDir(), "").Parse()
	if err == nil || !strings.Contains(err.Error(), "commands.json") || !strings.Contains(err.Error(), "--debug-output") {
		t.Fatalf("err: %v", err)
	}
}

func TestParse_SkipsFlowWithoutNameAndFileName(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "anon", header("", "", nil, nil, "U", "1"))
	writeFlow(t, root, "ok", okFlow("OK"))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].Title != "OK" {
		t.Fatalf("results: %+v", res)
	}
}

func TestParse_EndTimeIsLatestCommandEnd(t *testing.T) {
	root := t.TempDir()
	cmds := header("E", "f", nil, nil, "U", "1")
	cmds[0].Ts = 1000
	cmds[1].Ts = 1000
	cmds = append(cmds,
		cmdSpec{Type: "runFlowCommand", Status: "COMPLETED", Ts: 1000, Dur: 5000},
		cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Depth: 1, Ts: 2000, Dur: 100},
	)
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := res[0].Execution
	if *e.EndTime != 6000 || *e.Duration != 5000 {
		t.Errorf("end=%v dur=%v", *e.EndTime, *e.Duration)
	}
}

func TestParse_PendingCommandFailsFlow(t *testing.T) {
	root := t.TempDir()
	cmds := append(okFlow("P"), cmdSpec{Type: "tapOnElement", Status: "PENDING", Ts: 1300})
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Execution.Status != "failed" || res[0].Message == nil || *res[0].Message != "Flow did not finish" {
		t.Errorf("status %q message %v", res[0].Execution.Status, res[0].Message)
	}
}

func TestParse_NestedFailureWithCompletedParentPasses(t *testing.T) {
	root := t.TempDir()
	cmds := append(okFlow("N"),
		cmdSpec{Type: "runFlowCommand", Status: "COMPLETED", Ts: 1200, Dur: 50},
		cmdSpec{Type: "tapOnElement", Status: "FAILED", Depth: 1, Ts: 1210, Dur: 5},
	)
	writeFlow(t, root, "x", cmds)
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Execution.Status != "passed" {
		t.Errorf("status %q", res[0].Execution.Status)
	}
}
