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
