package client

import (
	"io"
	"net/http"
	"strings"

	"github.com/qase-tms/qasectl/internal"
)

// defaultHost is the Qase API host used when no host is configured.
const defaultHost = "api.qase.io"

// cliName is how the CLI identifies itself to the Qase API.
const cliName = "qasectl"

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

// userAgent builds the User-Agent sent to the Qase API. The CLI identifies
// itself first so that the backend can tell it apart from any other Go
// integration, and the generated client's own agent is kept as the suffix.
func userAgent(version, base string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		version = internal.UnknownVersion
	}

	agent := cliName + "/" + version

	base = strings.TrimSpace(base)
	if base == "" {
		return agent
	}

	return agent + " " + base
}

func extractBody(resp *http.Response) io.ReadCloser {
	if resp != nil {
		return resp.Body
	}
	return nil
}
