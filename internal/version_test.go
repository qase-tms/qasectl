package internal

import "testing"

func TestResolveVersionFrom(t *testing.T) {
	tests := []struct {
		name      string
		ldflags   string
		buildMain string
		want      string
	}{
		{
			name:    "build time version wins",
			ldflags: "v0.5.1",
			want:    "0.5.1",
		},
		{
			name:    "version without v prefix is kept as is",
			ldflags: "0.5.1",
			want:    "0.5.1",
		},
		{
			name:    "git describe output is preserved",
			ldflags: "v0.5.1-3-gabc1234",
			want:    "0.5.1-3-gabc1234",
		},
		{
			name:    "surrounding whitespace is trimmed",
			ldflags: "  v0.5.1  ",
			want:    "0.5.1",
		},
		{
			name:      "falls back to module version when not set at build time",
			ldflags:   "",
			buildMain: "v0.4.0",
			want:      "0.4.0",
		},
		{
			name:      "devel module version is not a version",
			ldflags:   "",
			buildMain: "(devel)",
			want:      "unknown",
		},
		{
			name:      "nothing known yields unknown",
			ldflags:   "",
			buildMain: "",
			want:      "unknown",
		},
		{
			name:      "build time version wins over module version",
			ldflags:   "v0.5.1",
			buildMain: "v0.4.0",
			want:      "0.5.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveVersionFrom(tt.ldflags, tt.buildMain)
			if got != tt.want {
				t.Errorf("resolveVersionFrom(%q, %q) = %q, want %q", tt.ldflags, tt.buildMain, got, tt.want)
			}
		})
	}
}

func TestResolveVersionIsNeverEmpty(t *testing.T) {
	if got := ResolveVersion(); got == "" {
		t.Error("ResolveVersion() = \"\", want a non-empty version")
	}
}
