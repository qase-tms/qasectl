package xctest

import (
	"fmt"
	"strings"
)

// stepAttachmentMode selects which steps keep the attachments XCTest records on activities.
// Attachments of failure summaries are always kept: they are the evidence of the failure.
type stepAttachmentMode int

const (
	FailedStepAttachments stepAttachmentMode = iota
	AllStepAttachments
	NoStepAttachments
)

func parseStepAttachmentMode(mode string) (stepAttachmentMode, error) {
	switch strings.ToLower(mode) {
	case "", "failed":
		return FailedStepAttachments, nil
	case "all":
		return AllStepAttachments, nil
	case "none":
		return NoStepAttachments, nil
	default:
		return 0, fmt.Errorf("unknown step attachments mode: %s. allowed values: failed, all, none", mode)
	}
}

func (m stepAttachmentMode) includes(isStepFailed bool) bool {
	switch m {
	case AllStepAttachments:
		return true
	case NoStepAttachments:
		return false
	default:
		return isStepFailed
	}
}
