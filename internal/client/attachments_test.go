package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

type attachmentServer struct {
	mu           sync.Mutex
	contentTypes map[string]string
	inFlight     atomic.Int32
	maxInFlight  atomic.Int32
	bulkBody     []byte
}

// newAttachmentServer answers attachment uploads with hash "h-<file name>" and records the
// Content-Type of every uploaded part.
func newAttachmentServer(t *testing.T) (*attachmentServer, *httptest.Server) {
	t.Helper()
	s := &attachmentServer{contentTypes: map[string]string{}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/attachment/") {
			n := s.inFlight.Add(1)
			defer s.inFlight.Add(-1)
			for {
				m := s.maxInFlight.Load()
				if n <= m || s.maxInFlight.CompareAndSwap(m, n) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)

			if err := r.ParseMultipartForm(1 << 20); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fh := r.MultipartForm.File["file[]"][0]
			s.mu.Lock()
			s.contentTypes[fh.Filename] = fh.Header.Get("Content-Type")
			s.mu.Unlock()

			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"status":true,"result":[{"hash":"h-%s","filename":%q}]}`, fh.Filename, fh.Filename)
			return
		}

		if strings.HasSuffix(r.URL.Path, "/results") {
			body := make([]byte, r.ContentLength)
			_, _ = r.Body.Read(body)
			s.mu.Lock()
			s.bulkBody = body
			s.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
			return
		}

		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	return s, srv
}

func content(s string) *[]byte {
	b := []byte(s)
	return &b
}

func TestClientV2_UploadData_SendsAttachmentHashes(t *testing.T) {
	s, srv := newAttachmentServer(t)
	c1 := NewClientV1("token", srv.URL)
	c2 := NewClientV2("token", srv.URL, c1)

	results := []models.Result{{
		Title:       "test",
		Execution:   models.Execution{Status: "failed"},
		Attachments: []models.Attachment{{Name: "result.txt", Content: content("r")}},
		Steps: []models.Step{{
			Data:      models.Data{Action: "parent"},
			Execution: models.StepExecution{Status: "failed"},
			Steps: []models.Step{{
				Data: models.Data{Action: "child"},
				Execution: models.StepExecution{
					Status:      "failed",
					Attachments: []models.Attachment{{Name: "child.plist", Content: content("c")}},
				},
			}},
		}},
	}}

	if err := c2.UploadData(context.Background(), "PRJ", 1, results); err != nil {
		t.Fatalf("UploadData() error = %v", err)
	}

	var body struct {
		Results []struct {
			Attachments []string `json:"attachments"`
			Steps       []struct {
				Steps []struct {
					Execution struct {
						Attachments []string `json:"attachments"`
					} `json:"execution"`
				} `json:"steps"`
			} `json:"steps"`
		} `json:"results"`
	}
	if err := json.Unmarshal(s.bulkBody, &body); err != nil {
		t.Fatalf("failed to decode bulk request %q: %v", s.bulkBody, err)
	}

	if got := body.Results[0].Attachments; len(got) != 1 || got[0] != "h-result.txt" {
		t.Errorf("result attachments = %v, want [h-result.txt]", got)
	}
	if got := body.Results[0].Steps[0].Steps[0].Execution.Attachments; len(got) != 1 || got[0] != "h-child.plist" {
		t.Errorf("nested step attachments = %v, want [h-child.plist]", got)
	}
	if got := s.contentTypes["child.plist"]; got != "application/x-plist" {
		t.Errorf("plist Content-Type = %q, want application/x-plist", got)
	}
}

func TestClientV1_uploadAttachments_RunsConcurrentlyWithinLimit(t *testing.T) {
	s, srv := newAttachmentServer(t)
	c := NewClientV1("token", srv.URL)

	var attachments []models.Attachment
	for i := range 3 * maxConcurrentAttachmentUploads {
		attachments = append(attachments, models.Attachment{Name: fmt.Sprintf("a%d.txt", i), Content: content("x")})
	}
	results := []models.Result{{Attachments: attachments}}

	c.uploadAttachments(context.Background(), "PRJ", results)

	for _, a := range results[0].Attachments {
		if a.Hash == nil || *a.Hash != "h-"+a.Name {
			t.Errorf("attachment %s hash = %v, want h-%s", a.Name, a.Hash, a.Name)
		}
	}
	if got := s.maxInFlight.Load(); got < 2 || got > maxConcurrentAttachmentUploads {
		t.Errorf("max concurrent uploads = %d, want between 2 and %d", got, maxConcurrentAttachmentUploads)
	}
}

func TestClientV1_uploadAttachments_SameNameInDifferentResults(t *testing.T) {
	s, srv := newAttachmentServer(t)
	c := NewClientV1("token", srv.URL)

	results := []models.Result{
		{Attachments: []models.Attachment{{Name: "screenshot.jpeg", Content: content("1")}}},
		{Attachments: []models.Attachment{{Name: "screenshot.jpeg", Content: content("2")}}},
	}

	c.uploadAttachments(context.Background(), "PRJ", results)

	for i, r := range results {
		if r.Attachments[0].Hash == nil {
			t.Errorf("result %d attachment was not uploaded", i)
		}
	}
	if got := s.contentTypes["screenshot.jpeg"]; got != "image/jpeg" {
		t.Errorf("Content-Type = %q, want image/jpeg", got)
	}
}

func TestClientV1_uploadAttachments_KeepsFileOnDisk(t *testing.T) {
	_, srv := newAttachmentServer(t)
	c := NewClientV1("token", srv.URL)

	path := filepath.Join(t.TempDir(), "report.pdf")
	if err := os.WriteFile(path, []byte("%PDF"), 0o600); err != nil {
		t.Fatal(err)
	}
	results := []models.Result{{Attachments: []models.Attachment{{Name: "report.pdf", FilePath: &path}}}}

	c.uploadAttachments(context.Background(), "PRJ", results)

	if h := results[0].Attachments[0].Hash; h == nil || *h != "h-report.pdf" {
		t.Errorf("hash = %v, want h-report.pdf", h)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("attachment file given by path was removed: %v", err)
	}
}

func TestClientV1_uploadAttachments_SkipsFailedUploads(t *testing.T) {
	_, srv := newAttachmentServer(t)
	c := NewClientV1("token", srv.URL)

	missing := filepath.Join(t.TempDir(), "missing.txt")
	results := []models.Result{{Attachments: []models.Attachment{
		{Name: "missing.txt", FilePath: &missing},
		{Name: "empty"},
		{Name: "ok.txt", Content: content("ok")},
	}}}

	c.uploadAttachments(context.Background(), "PRJ", results)

	got := c.convertAttachments(context.Background(), "PRJ", results[0].Attachments)
	if len(got) != 1 || got[0] != "h-ok.txt" {
		t.Errorf("convertAttachments() = %v, want [h-ok.txt]", got)
	}
}

func TestAttachmentFileName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "screenshot.jpeg", want: "screenshot.jpeg"},
		{name: "a/b\\c.txt", want: "a_b_c.txt"},
		{name: "", want: "attachment"},
		{name: "..", want: "attachment"},
	}
	for _, tt := range tests {
		if got := attachmentFileName(tt.name); got != tt.want {
			t.Errorf("attachmentFileName(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
