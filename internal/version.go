package internal

import (
	"runtime/debug"
	"strings"
)

// UnknownVersion is reported when the version cannot be determined.
const UnknownVersion = "unknown"

// develVersion is what the toolchain reports for a module built from a
// working tree rather than from a tagged release.
const develVersion = "(devel)"

// ResolveVersion returns the CLI version.
// It prefers the version stamped at build time via ldflags and falls back to
// the module version recorded in the binary, which covers "go install".
func ResolveVersion() string {
	var buildMain string
	if info, ok := debug.ReadBuildInfo(); ok {
		buildMain = info.Main.Version
	}

	return resolveVersionFrom(Version, buildMain)
}

// resolveVersionFrom picks the first usable version and normalises it.
func resolveVersionFrom(ldflags, buildMain string) string {
	for _, candidate := range []string{ldflags, buildMain} {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" || candidate == develVersion {
			continue
		}

		return strings.TrimPrefix(candidate, "v")
	}

	return UnknownVersion
}
