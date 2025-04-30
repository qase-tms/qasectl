package client

import (
	"io"
	"net/http"
)

func extractBody(resp *http.Response) io.ReadCloser {
	if resp != nil {
		return resp.Body
	}
	return nil
}
