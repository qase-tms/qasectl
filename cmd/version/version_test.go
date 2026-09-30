package version

import (
	"testing"

	"github.com/qase-tms/qasectl/internal"
)

func TestVersionOutput(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "known version is printed",
			version: "0.5.1",
			want:    "Qase CLI 0.5.1",
		},
		{
			name:    "unknown version explains itself",
			version: internal.UnknownVersion,
			want:    "Qase CLI: version not set during build.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := versionOutput(tt.version)
			if got != tt.want {
				t.Errorf("versionOutput(%q) = %q, want %q", tt.version, got, tt.want)
			}
		})
	}
}
