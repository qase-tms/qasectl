package client

import "testing"

func TestApiURL(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		version string
		want    string
	}{
		{
			name:    "empty host falls back to default",
			host:    "",
			version: "v1",
			want:    "https://api.qase.io/v1",
		},
		{
			name:    "default host v2",
			host:    "",
			version: "v2",
			want:    "https://api.qase.io/v2",
		},
		{
			name:    "plain host gets https scheme",
			host:    "api.qase.io",
			version: "v1",
			want:    "https://api.qase.io/v1",
		},
		{
			name:    "custom enterprise host",
			host:    "api.qase.example.com",
			version: "v2",
			want:    "https://api.qase.example.com/v2",
		},
		{
			name:    "trailing slash is trimmed",
			host:    "api.qase.example.com/",
			version: "v1",
			want:    "https://api.qase.example.com/v1",
		},
		{
			name:    "explicit https scheme is preserved",
			host:    "https://api.qase.example.com",
			version: "v1",
			want:    "https://api.qase.example.com/v1",
		},
		{
			name:    "explicit http scheme is preserved",
			host:    "http://localhost:8080",
			version: "v2",
			want:    "http://localhost:8080/v2",
		},
		{
			name:    "surrounding whitespace is trimmed",
			host:    "  api.qase.io  ",
			version: "v1",
			want:    "https://api.qase.io/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiURL(tt.host, tt.version)
			if got != tt.want {
				t.Errorf("apiURL(%q, %q) = %q, want %q", tt.host, tt.version, got, tt.want)
			}
		})
	}
}

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name    string
		version string
		base    string
		want    string
	}{
		{
			name:    "cli identifies itself ahead of the generated client",
			version: "0.5.1",
			base:    "qase-api-client-go/1.2.13",
			want:    "qasectl/0.5.1 qase-api-client-go/1.2.13",
		},
		{
			name:    "empty base leaves no trailing space",
			version: "0.5.1",
			base:    "",
			want:    "qasectl/0.5.1",
		},
		{
			name:    "missing version is reported as unknown",
			version: "",
			base:    "qase-api-client-go/1.2.13",
			want:    "qasectl/unknown qase-api-client-go/1.2.13",
		},
		{
			name:    "surrounding whitespace is trimmed",
			version: "0.5.1",
			base:    "  qase-api-client-go/1.2.13  ",
			want:    "qasectl/0.5.1 qase-api-client-go/1.2.13",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := userAgent(tt.version, tt.base)
			if got != tt.want {
				t.Errorf("userAgent(%q, %q) = %q, want %q", tt.version, tt.base, got, tt.want)
			}
		})
	}
}
