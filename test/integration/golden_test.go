package integration

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden files")

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func normalize(reqs []recordedRequest) []byte {
	raw, err := json.Marshal(reqs)
	if err != nil {
		panic(err)
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(walk(v)); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func walk(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			x[k] = walk(e)
		}
		return x
	case []any:
		for i, e := range x {
			x[i] = walk(e)
		}
		return x
	case string:
		return normalizeString(x)
	}
	return v
}

func normalizeString(s string) string {
	if uuidRe.MatchString(s) {
		return "<uuid>"
	}
	roots := []string{repoRoot}
	if home, err := os.UserHomeDir(); err == nil {
		roots = append(roots, home)
	}
	return normalizeStringRoots(s, roots)
}

func normalizeStringRoots(s string, roots []string) string {
	for _, r := range roots {
		if r == "" || r == "/" {
			continue
		}
		if s == r || strings.HasPrefix(s, r+"/") {
			return "/workspace" + filepath.ToSlash(strings.TrimPrefix(s, r))
		}
	}
	return s
}

func assertGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	assertGoldenTB(t, path, got)
}

func assertGoldenTB(t testing.TB, path string, got []byte) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
			return
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Fatalf("golden %s missing: run go test ./test/integration/ -update", path)
		return
	}
	if err != nil {
		t.Fatalf("read golden: %v", err)
		return
	}
	if !bytes.Equal(want, got) {
		t.Errorf("golden %s mismatch:\n%s", path, lineDiff(string(want), string(got)))
	}
}

func lineDiff(want, got string) string {
	w := strings.Split(want, "\n")
	g := strings.Split(got, "\n")
	var out []string
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		wok, gok := i < len(w), i < len(g)
		if wok {
			wl = w[i]
		}
		if gok {
			gl = g[i]
		}
		if wok && gok && wl == gl {
			continue
		}
		if wok {
			out = append(out, fmt.Sprintf("%d: -%s", i+1, wl))
		}
		if gok {
			out = append(out, fmt.Sprintf("%d: +%s", i+1, gl))
		}
		if len(out) >= 40 {
			out = append(out[:40], "...")
			break
		}
	}
	return strings.Join(out, "\n")
}
