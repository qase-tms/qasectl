package maestro

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

const maxAttachmentSize = 32 << 20

func stepAttachments(flowDir string, c command) []models.Attachment {
	var out []models.Attachment
	for _, a := range c.Artifacts {
		if a.Type != "SCREENSHOT" && a.Type != "TAKE_SCREENSHOT" {
			continue
		}
		p, ok := insideFlow(flowDir, a.Path)
		if !ok {
			continue
		}
		if att, ok := attachFile(p); ok {
			out = append(out, att)
		}
	}
	return out
}

func failureLogs(flowDir string) []models.Attachment {
	paths := []string{filepath.Join(flowDir, "logs", "maestro.log")}
	devices, _ := filepath.Glob(filepath.Join(flowDir, "logs", "device-*.log"))
	paths = append(paths, devices...)

	var out []models.Attachment
	for _, p := range paths {
		if att, ok := attachFile(p); ok {
			out = append(out, att)
		}
	}
	return out
}

func insideFlow(flowDir, rel string) (string, bool) {
	if filepath.IsAbs(rel) {
		slog.Warn("skipping Maestro artifact with absolute path", "path", rel)
		return "", false
	}
	p := filepath.Join(flowDir, rel)
	r, err := filepath.Rel(flowDir, p)
	if err != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		slog.Warn("skipping Maestro artifact outside flow directory", "path", rel)
		return "", false
	}
	return p, true
}

func attachFile(path string) (models.Attachment, bool) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		slog.Warn("Maestro attachment not found, skipping", "path", path)
		return models.Attachment{}, false
	}
	name := filepath.Base(path)
	if info.Size() <= maxAttachmentSize {
		p := path
		return models.Attachment{Name: name, FilePath: &p, ContentType: contentType(name)}, true
	}

	data, err := gzipFile(path)
	if err != nil {
		slog.Warn("failed to compress Maestro attachment, skipping", "path", path, "error", err)
		return models.Attachment{}, false
	}
	if len(data) > maxAttachmentSize {
		slog.Warn("Maestro attachment exceeds 32 MB even after gzip, skipping", "path", path)
		return models.Attachment{}, false
	}
	return models.Attachment{Name: name + ".gz", Content: &data, ContentType: "application/gzip"}, true
}

func gzipFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := io.Copy(zw, f); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func contentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".log", ".txt":
		return "text/plain"
	}
	return "application/octet-stream"
}
