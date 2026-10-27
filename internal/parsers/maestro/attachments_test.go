package maestro

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func arts(pairs ...string) command {
	c := command{}
	for i := 0; i < len(pairs); i += 2 {
		c.Artifacts = append(c.Artifacts, artifact{Type: pairs[i], Path: pairs[i+1]})
	}
	return c
}

func TestStepAttachments_Screenshots(t *testing.T) {
	flowDir := t.TempDir()
	writeFile(t, flowDir, "screenshots/a.png", 10)
	writeFile(t, flowDir, "screenshots/b.png", 10)
	writeFile(t, flowDir, "screenshots/h.json", 10)

	got := stepAttachments(flowDir, arts(
		"SCREENSHOT", "screenshots/a.png",
		"TAKE_SCREENSHOT", "screenshots/b.png",
		"SCREEN_HIERARCHY", "screenshots/h.json",
	))
	if len(got) != 2 {
		t.Fatalf("got %d attachments", len(got))
	}
	for _, a := range got {
		if a.FilePath == nil || !strings.HasPrefix(*a.FilePath, flowDir) {
			t.Errorf("file path: %v", a.FilePath)
		}
		if a.ContentType != "image/png" {
			t.Errorf("content type: %s", a.ContentType)
		}
	}
}

func TestStepAttachments_MissingFileSkipped(t *testing.T) {
	if got := stepAttachments(t.TempDir(), arts("SCREENSHOT", "screenshots/none.png")); len(got) != 0 {
		t.Errorf("got %d attachments", len(got))
	}
}

func TestStepAttachments_RejectsPathOutsideFlow(t *testing.T) {
	root := t.TempDir()
	flowDir := filepath.Join(root, "flow")
	if err := os.MkdirAll(flowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	secret := writeFile(t, root, "secret.png", 10)

	got := stepAttachments(flowDir, arts("SCREENSHOT", "../secret.png", "SCREENSHOT", secret))
	if len(got) != 0 {
		t.Errorf("got %d attachments", len(got))
	}
}

func TestFailureLogs_OnlyForFailedFlows(t *testing.T) {
	root := t.TempDir()
	fail := append(okFlow("Failing"), cmdSpec{Type: "assertVisible", Status: "FAILED", Ts: 1200})
	failDir := writeFlow(t, root, "fail", fail)
	passDir := writeFlow(t, root, "pass", okFlow("Passing"))
	for _, d := range []string{failDir, passDir} {
		for _, l := range []string{"maestro.log", "device-simulator.log", "device-xctest.log", "other.log"} {
			writeFile(t, d, "logs/"+l, 5)
		}
	}
	writeFile(t, root, "maestro.log", 5)

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		want := 0
		if r.Title == "Failing" {
			want = 3
		}
		if len(r.Attachments) != want {
			t.Errorf("%s: got %d attachments, want %d", r.Title, len(r.Attachments), want)
		}
	}
}

func TestAttachFile_GzipOverLimit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "big.log")
	if err := os.WriteFile(p, bytes.Repeat([]byte("a"), maxAttachmentSize+1), 0o644); err != nil {
		t.Fatal(err)
	}

	a, ok := attachFile(p)
	if !ok {
		t.Fatal("not attached")
	}
	if a.Name != "big.log.gz" || a.ContentType != "application/gzip" || a.FilePath != nil {
		t.Errorf("unexpected attachment: %+v", a)
	}
	if a.Content == nil || len(*a.Content) >= maxAttachmentSize {
		t.Fatalf("content not compressed")
	}
	zr, err := gzip.NewReader(bytes.NewReader(*a.Content))
	if err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != maxAttachmentSize+1 || out[0] != 'a' || out[len(out)-1] != 'a' {
		t.Errorf("round trip mismatch: %d bytes", len(out))
	}
}

func TestAttachFile_UnderLimitUsesPath(t *testing.T) {
	p := writeFile(t, t.TempDir(), "x.log", 10)
	a, ok := attachFile(p)
	if !ok {
		t.Fatal("not attached")
	}
	if a.FilePath == nil || *a.FilePath != p || a.Content != nil || a.Name != "x.log" || a.ContentType != "text/plain" {
		t.Errorf("unexpected attachment: %+v", a)
	}
}

func TestParse_ScreenshotOnStep(t *testing.T) {
	root := t.TempDir()
	cmds := header("Shot", "f", nil, nil, "UDID", "1")
	cmds = append(cmds, cmdSpec{
		Type: "assertVisible", Status: "WARNED", Ts: 1100,
		Artifacts: []map[string]string{{"type": "SCREENSHOT", "path": "screenshots/s.png"}},
	})
	dir := writeFlow(t, root, "shot", cmds)
	writeFile(t, dir, "screenshots/s.png", 4)

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if got := res[0].Steps[0].Execution.Attachments; len(got) != 1 {
		t.Errorf("step attachments: %d", len(got))
	}
	if len(res[0].Attachments) != 0 {
		t.Errorf("result attachments: %d", len(res[0].Attachments))
	}
}
