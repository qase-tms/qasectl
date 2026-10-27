package integration

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type captureTB struct {
	testing.TB
	msgs []string
}

func (c *captureTB) Errorf(format string, args ...any) {
	c.msgs = append(c.msgs, fmt.Sprintf(format, args...))
}
func (c *captureTB) Fatalf(format string, args ...any) {
	c.msgs = append(c.msgs, fmt.Sprintf(format, args...))
}

func postFile(t *testing.T, url, name, content string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file[]"; filename="%s"`, name))
	h.Set("Content-Type", "image/png")
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte(content))
	_ = w.Close()
	resp, err := http.Post(url+"/v1/attachment/TEST", w.FormDataContentType(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestFakeQase_UnexpectedRequestFails(t *testing.T) {
	ctb := &captureTB{TB: t}
	url, _ := startFake(ctb)

	resp, err := http.Get(url + "/v1/project/TEST")
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
	if len(ctb.msgs) != 1 || !strings.Contains(ctb.msgs[0], "GET /v1/project/TEST") {
		t.Errorf("errors = %v, want one mentioning GET /v1/project/TEST", ctb.msgs)
	}
}

func TestRecorder_SortsAttachmentUploads(t *testing.T) {
	url, rec := newFakeQase(t)
	postFile(t, url, "b.png", "bbb")
	postFile(t, url, "a.png", "aaa")

	reqs := rec.Requests()
	if len(reqs) != 2 {
		t.Fatalf("got %d requests, want 2", len(reqs))
	}
	if got := reqs[0].Files[0].Name; got != "a.png" {
		t.Errorf("first upload = %s, want a.png", got)
	}
	if reqs[0].Files[0].Size != 3 || reqs[0].Files[0].Ext != ".png" {
		t.Errorf("unexpected file record %+v", reqs[0].Files[0])
	}
}

func TestNormalize_ReplacesAbsolutePaths(t *testing.T) {
	reqs := []recordedRequest{{
		Method: "POST",
		Path:   "/v2/TEST/run/1/results",
		Body: map[string]any{
			"id":   "8a1b2c3d-1111-4222-8333-444455556666",
			"file": repoRoot + "/x/y.py",
		},
	}}
	want := `[
  {
    "body": {
      "file": "/workspace/x/y.py",
      "id": "<uuid>"
    },
    "method": "POST",
    "path": "/v2/TEST/run/1/results"
  }
]
`
	if got := string(normalize(reqs)); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestAssertGolden_MissingFile(t *testing.T) {
	ctb := &captureTB{TB: t}
	assertGoldenTB(ctb, filepath.Join(t.TempDir(), "nope.json"), []byte("{}\n"))
	if len(ctb.msgs) != 1 || !strings.Contains(ctb.msgs[0], "-update") {
		t.Errorf("messages = %v, want one mentioning -update", ctb.msgs)
	}
}

func TestAssertGolden_MismatchShowsDiff(t *testing.T) {
	p := filepath.Join(t.TempDir(), "g.json")
	if err := os.WriteFile(p, []byte("a\nb\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctb := &captureTB{TB: t}
	assertGoldenTB(ctb, p, []byte("a\nc\n"))
	if len(ctb.msgs) != 1 || !strings.Contains(ctb.msgs[0], "-b") || !strings.Contains(ctb.msgs[0], "+c") {
		t.Errorf("messages = %v", ctb.msgs)
	}
}

func TestRunUpload_FailurePrintsOutput(t *testing.T) {
	ctb := &captureTB{TB: t}
	runUploadTB(ctb, "http://127.0.0.1:1", "junit", "/nonexistent")
	if len(ctb.msgs) != 1 || !strings.Contains(ctb.msgs[0], "failed to parse results") {
		t.Errorf("messages = %v, want binary error text", ctb.msgs)
	}
}

func TestRecorder_SortIsDeterministicForSameName(t *testing.T) {
	var got [2][]recordedRequest
	for i, order := range [][]string{{"aa", "bbbb"}, {"bbbb", "aa"}} {
		url, rec := newFakeQase(t)
		for _, c := range order {
			postFile(t, url, "same.png", c)
		}
		got[i] = rec.Requests()
	}
	if string(normalize(got[0])) != string(normalize(got[1])) {
		t.Errorf("recordings differ by arrival order:\n%s\nvs\n%s", normalize(got[0]), normalize(got[1]))
	}
}

func TestNormalizeString_RootSlashIsSkipped(t *testing.T) {
	const in = "/workspace/x/y.py"
	if got := normalizeStringRoots(in, []string{"/", ""}); got != in {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestNormalizeString_MatchesOnPathBoundary(t *testing.T) {
	roots := []string{"/a/qasectl"}
	if got := normalizeStringRoots("/a/qasectl-other/x", roots); got != "/a/qasectl-other/x" {
		t.Errorf("sibling dir rewritten: %q", got)
	}
	if got := normalizeStringRoots("/a/qasectl/x", roots); got != "/workspace/x" {
		t.Errorf("got %q, want /workspace/x", got)
	}
	if got := normalizeStringRoots("/a/qasectl", roots); got != "/workspace" {
		t.Errorf("got %q, want /workspace", got)
	}
}

func TestAssertGolden_UpdateCreatesMissingFile(t *testing.T) {
	old := *update
	*update = true
	defer func() { *update = old }()

	p := filepath.Join(t.TempDir(), "nested", "g.json")
	assertGoldenTB(t, p, []byte("{}\n"))
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("golden not created: %v", err)
	}
	if string(got) != "{}\n" {
		t.Errorf("content = %q", got)
	}
}
