package integration

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

type recordedFile struct {
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

type recordedRequest struct {
	Method string         `json:"method"`
	Path   string         `json:"path"`
	Body   any            `json:"body,omitempty"`
	Files  []recordedFile `json:"files,omitempty"`
}

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

func (r *recorder) add(req recordedRequest) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reqs = append(r.reqs, req)
}

// Requests returns requests in arrival order, except attachment uploads,
// which are sorted by file name because the client uploads concurrently.
func (r *recorder) Requests() []recordedRequest {
	r.mu.Lock()
	out := append([]recordedRequest(nil), r.reqs...)
	r.mu.Unlock()

	isUpload := func(q recordedRequest) bool { return len(q.Files) > 0 }
	var slots []int
	var uploads []recordedRequest
	for i, q := range out {
		if isUpload(q) {
			slots = append(slots, i)
			uploads = append(uploads, q)
		}
	}
	sort.SliceStable(uploads, func(i, j int) bool { return filesLess(uploads[i].Files, uploads[j].Files) })
	for n, i := range slots {
		out[i] = uploads[n]
	}
	return out
}

func filesLess(a, b []recordedFile) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		x, y := a[i], b[i]
		if x.Name != y.Name {
			return x.Name < y.Name
		}
		if x.Ext != y.Ext {
			return x.Ext < y.Ext
		}
		if x.Size != y.Size {
			return x.Size < y.Size
		}
	}
	return len(a) < len(b)
}

func newFakeQase(t *testing.T) (string, *recorder) {
	t.Helper()
	return startFake(t)
}

func startFake(t testing.TB) (string, *recorder) {
	t.Helper()
	rec := &recorder{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.Method + " " + r.URL.Path
		switch route {
		case "POST /v1/run/TEST":
			rec.add(recordedRequest{Method: r.Method, Path: r.URL.Path, Body: readJSON(r)})
			writeJSON(w, `{"status":true,"result":{"id":1}}`)
		case "POST /v1/attachment/TEST":
			handleAttachments(w, r, rec)
		case "POST /v2/TEST/run/1/results":
			rec.add(recordedRequest{Method: r.Method, Path: r.URL.Path, Body: readJSON(r)})
			writeJSON(w, `{"status":true}`)
		case "POST /v1/run/TEST/1/complete":
			rec.add(recordedRequest{Method: r.Method, Path: r.URL.Path})
			writeJSON(w, `{"status":true}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL, rec
}

func readJSON(r *http.Request) any {
	b, _ := io.ReadAll(r.Body)
	if len(b) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	return v
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}

func handleAttachments(w http.ResponseWriter, r *http.Request, rec *recorder) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var files []recordedFile
	var results []map[string]any
	for _, fh := range r.MultipartForm.File["file[]"] {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		content, _ := io.ReadAll(f)
		_ = f.Close()

		name := filepath.Base(fh.Filename)
		mime := fh.Header.Get("Content-Type")
		ext := strings.ToLower(filepath.Ext(name))
		sum := sha1.Sum(append([]byte(name), content...))
		hash := hex.EncodeToString(sum[:])
		files = append(files, recordedFile{Name: name, Ext: ext, Size: int64(len(content))})
		results = append(results, map[string]any{
			"hash": hash, "filename": name, "mime": mime, "size": len(content),
			"extension": ext, "url": "http://fake/" + hash,
		})
	}
	rec.add(recordedRequest{Method: r.Method, Path: r.URL.Path, Files: files})
	b, _ := json.Marshal(map[string]any{"status": true, "result": results})
	writeJSON(w, string(b))
}
