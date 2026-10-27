package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	qasectlBin string
	repoRoot   string
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "go env GOMOD:", err)
		return 1
	}
	repoRoot = filepath.Dir(strings.TrimSpace(string(out)))

	tmp, err := os.MkdirTemp("", "qasectl-bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	qasectlBin = filepath.Join(tmp, "qasectl")
	build := exec.Command("go", "build", "-o", qasectlBin, ".")
	build.Dir = repoRoot
	if b, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s", err, b)
		return 1
	}
	return m.Run()
}

func runUpload(t *testing.T, apiHost, format, path string, extra ...string) string {
	t.Helper()
	return runUploadTB(t, apiHost, format, path, extra...)
}

func runUploadTB(t testing.TB, apiHost, format, path string, extra ...string) string {
	t.Helper()
	args := append([]string{
		"testops", "result", "upload",
		"--project", "TEST", "--token", "test", "--api-host", apiHost,
		"--title", "Integration", "--format", format, "--path", path,
	}, extra...)
	cmd := exec.Command(qasectlBin, args...)
	cmd.Env = cleanEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("qasectl %s: %v\n%s", strings.Join(args, " "), err, out)
		return string(out)
	}
	return string(out)
}

func cleanEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "QASE_TESTOPS_") {
			env = append(env, kv)
		}
	}
	return env
}
