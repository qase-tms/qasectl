package xctest

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// xcresulttoolRequiringLegacy answers like xcresulttool on Xcode 16+: a legacy command
// exits with status 64 unless --legacy is passed.
const xcresulttoolRequiringLegacy = `case " $* " in
*" --legacy "*) echo '{}' ;;
*) echo "Error: This command is deprecated and will be removed in a future release, --legacy flag is required to use it." >&2; exit 64 ;;
esac
`

// stubXcrun points the parser at a shell script standing in for xcrun. The script records
// each invocation's arguments and then runs body; the returned function reads them back.
func stubXcrun(t *testing.T, p *Parser, body string) func() []string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the xcrun stub is a POSIX shell script")
	}

	dir := t.TempDir()
	calls := filepath.Join(dir, "calls")
	p.xcrun = filepath.Join(dir, "xcrun")
	script := "#!/bin/sh\necho \"$*\" >> \"" + calls + "\"\n" + body
	if err := os.WriteFile(p.xcrun, []byte(script), 0o755); err != nil {
		t.Fatalf("Failed to write xcrun stub: %v", err)
	}

	return func() []string {
		data, err := os.ReadFile(calls)
		if err != nil {
			t.Fatalf("Failed to read recorded xcrun calls: %v", err)
		}
		return strings.Split(strings.TrimSpace(string(data)), "\n")
	}
}

// TestParser_xcresulttool_RemembersLegacyFlag verifies that only the first call pays for
// the invocation xcresulttool rejects without --legacy.
func TestParser_xcresulttool_RemembersLegacyFlag(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	calls := stubXcrun(t, p, xcresulttoolRequiringLegacy)

	id := "0~test"
	for range 3 {
		if _, err := p.readJson(&id); err != nil {
			t.Fatalf("readJson() error = %v", err)
		}
	}
	if _, err := p.readAttachment("0~attachment"); err != nil {
		t.Fatalf("readAttachment() error = %v", err)
	}

	got := calls()
	if len(got) != 5 {
		t.Fatalf("xcrun ran %d times, want 5 (one rejected call, then four with --legacy): %q", len(got), got)
	}
	if strings.Contains(got[0], "--legacy") {
		t.Errorf("first call = %q, want it without --legacy", got[0])
	}
	for _, call := range got[1:] {
		if !strings.HasSuffix(call, " --legacy") {
			t.Errorf("call = %q, want it with --legacy", call)
		}
	}
}

// TestParser_xcresulttool_OmitsLegacyWhenNotRequired verifies that an xcresulttool which
// accepts plain calls (Xcode 15 and earlier) never receives --legacy.
func TestParser_xcresulttool_OmitsLegacyWhenNotRequired(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	calls := stubXcrun(t, p, "echo '{}'\n")

	for range 2 {
		if _, err := p.readJson(nil); err != nil {
			t.Fatalf("readJson() error = %v", err)
		}
	}

	got := calls()
	if len(got) != 2 {
		t.Fatalf("xcrun ran %d times, want 2: %q", len(got), got)
	}
	for _, call := range got {
		if strings.Contains(call, "--legacy") {
			t.Errorf("call = %q, want it without --legacy", call)
		}
	}
}

// TestParser_xcresulttool_DoesNotRetryOtherFailures verifies that only exit status 64
// triggers the --legacy retry.
func TestParser_xcresulttool_DoesNotRetryOtherFailures(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	calls := stubXcrun(t, p, "echo 'Error: unexpected failure' >&2\nexit 1\n")

	if _, err := p.readJson(nil); err == nil {
		t.Fatal("readJson() error = nil, want the failure")
	}

	if got := calls(); len(got) != 1 {
		t.Errorf("xcrun ran %d times, want 1: %q", len(got), got)
	}
}

// TestParser_xcresulttool_ErrorIncludesStderr verifies that a failed --legacy retry reports
// xcresulttool's message rather than a bare exit status, and is not remembered.
func TestParser_xcresulttool_ErrorIncludesStderr(t *testing.T) {
	p, err := NewParser("test.xcresult", "")
	if err != nil {
		t.Fatalf("Failed to create parser: %v", err)
	}
	stubXcrun(t, p, `case " $* " in
*" --legacy "*) echo "Error: File or directory doesn't exist at path: test.xcresult." >&2; exit 64 ;;
*) echo "Error: --legacy flag is required to use it." >&2; exit 64 ;;
esac
`)

	_, err = p.readJson(nil)
	if err == nil {
		t.Fatal("readJson() error = nil, want the missing-path failure")
	}
	if want := "File or directory doesn't exist at path: test.xcresult."; !strings.Contains(err.Error(), want) {
		t.Errorf("readJson() error = %q, want it to contain %q", err, want)
	}
	if p.legacy {
		t.Error("legacy = true after a failed retry, want false")
	}
}
