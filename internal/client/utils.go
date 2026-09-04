package client

import (
	"io"
	"net/http"
	"strings"
)

// defaultHost is the Qase API host used when no host is configured.
const defaultHost = "api.qase.io"

// apiURL builds a base URL for the given Qase API version.
// The host may be a bare host name ("api.qase.io"), in which case the https
// scheme is added, or carry an explicit scheme ("http://localhost:8080"),
// which is preserved. An empty host falls back to defaultHost.
func apiURL(host, version string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		host = defaultHost
	}

	host = strings.TrimRight(host, "/")
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}

	return host + "/" + version
}

func extractBody(resp *http.Response) io.ReadCloser {
	if resp != nil {
		return resp.Body
	}
	return nil
}
