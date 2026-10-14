package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

const maxConcurrentAttachmentUploads = 8

func init() {
	// The API client derives each file's Content-Type from its extension, and Go knows no
	// type for the property lists XCTest attaches.
	_ = mime.AddExtensionType(".plist", "application/x-plist")
}

// uploadAttachments uploads every attachment of results, including step attachments,
// and stores the returned hash on the attachment.
func (c *ClientV1) uploadAttachments(ctx context.Context, project string, results []models.Result) {
	const op = "client.attachments.uploadattachments"
	logger := slog.With("op", op)

	attachments := collectAttachments(results)
	if len(attachments) == 0 {
		return
	}

	logger.Info("uploading attachments", "count", len(attachments))

	var wg sync.WaitGroup
	for _, a := range attachments {
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case c.uploads <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-c.uploads }()

			hash, err := c.uploadAttachmentData(ctx, project, *a)
			if err != nil {
				logger.Warn("failed to upload attachment", "name", a.Name, "error", err)
				return
			}
			a.Hash = &hash
		}()
	}
	wg.Wait()
}

func collectAttachments(results []models.Result) []*models.Attachment {
	var attachments []*models.Attachment
	var collectSteps func(steps []models.Step)
	collectSteps = func(steps []models.Step) {
		for i := range steps {
			for j := range steps[i].Execution.Attachments {
				attachments = append(attachments, &steps[i].Execution.Attachments[j])
			}
			collectSteps(steps[i].Steps)
		}
	}

	for i := range results {
		for j := range results[i].Attachments {
			attachments = append(attachments, &results[i].Attachments[j])
		}
		collectSteps(results[i].Steps)
	}

	return attachments
}

// uploadAttachmentData writes in-memory content to a private temp dir: the uploaded file
// name is the attachment name, and attachments of different results may share it.
func (c *ClientV1) uploadAttachmentData(ctx context.Context, project string, a models.Attachment) (string, error) {
	path := a.FilePath
	if path == nil {
		if a.Content == nil {
			return "", errors.New("attachment has neither a file path nor content")
		}

		dir, err := os.MkdirTemp("", "qasectl-attachment-*")
		if err != nil {
			return "", fmt.Errorf("failed to create temp dir: %w", err)
		}
		defer func() { _ = os.RemoveAll(dir) }()

		fp := filepath.Join(dir, attachmentFileName(a.Name))
		if err := os.WriteFile(fp, *a.Content, 0o600); err != nil {
			return "", fmt.Errorf("failed to write attachment: %w", err)
		}
		path = &fp
	}

	file, err := os.Open(*path)
	if err != nil {
		return "", fmt.Errorf("failed to open attachment: %w", err)
	}
	defer func() { _ = file.Close() }()

	return c.uploadAttachment(ctx, project, []*os.File{file})
}

func attachmentFileName(name string) string {
	name = strings.NewReplacer("/", "_", "\\", "_").Replace(name)
	if name == "" || name == "." || name == ".." {
		return "attachment"
	}

	return name
}
