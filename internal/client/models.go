package client

import (
	"fmt"
	"io"
	"log"
)

// QaseApiError is an error returned by Qase API
type QaseApiError struct {
	message string
	reason  string
}

// NewQaseApiError creates a new QaseApiError
func NewQaseApiError(message string, reason io.ReadCloser) *QaseApiError {
	bodyBytes := []byte("")
	if reason != nil {
		var err error
		bodyBytes, err = io.ReadAll(reason)
		if err != nil {
			log.Fatalf("Failed to read body: %v", err)
		}
	}

	return &QaseApiError{
		message: message,
		reason:  string(bodyBytes),
	}
}

// Error returns a string representation of the error
func (e *QaseApiError) Error() string {
	return fmt.Sprintf("Message: %s. Reason: %s.", e.message, e.reason)
}
