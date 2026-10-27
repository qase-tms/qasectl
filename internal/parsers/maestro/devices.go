package maestro

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const deviceSep = " - "

type junitSuites struct {
	Suites []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	Device string `xml:"device,attr"`
}

// deviceNames maps UDID to "<model> - <OS>" from Maestro JUnit reports under root.
func deviceNames(root string) map[string]string {
	names := map[string]string{}
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.EqualFold(filepath.Ext(d.Name()), ".xml") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, s := range readSuites(data) {
			if i := strings.LastIndex(s.Device, deviceSep); i > 0 {
				if udid := strings.TrimSpace(s.Device[i+len(deviceSep):]); udid != "" {
					names[udid] = s.Device[:i]
				}
			}
		}
		return nil
	})
	return names
}

func readSuites(data []byte) []junitSuite {
	var multi junitSuites
	if xml.Unmarshal(data, &multi) == nil && len(multi.Suites) > 0 {
		return multi.Suites
	}
	var single junitSuite
	if xml.Unmarshal(data, &single) == nil && single.Device != "" {
		return []junitSuite{single}
	}
	return nil
}

func checkDuplicates(flows []flow) error {
	type key struct{ name, udid string }
	seen := map[key]bool{}
	for _, f := range flows {
		name := f.Name
		if name == "" {
			name = f.FileName
		}
		k := key{name, f.UDID}
		if seen[k] {
			return errors.New("path contains several Maestro runs; pass the folder of a single run")
		}
		seen[k] = true
	}
	return nil
}

func hasSeveralShards(flows []flow) bool {
	ids := map[string]bool{}
	for _, f := range flows {
		if f.ShardID != "" {
			ids[f.ShardID] = true
		}
	}
	return len(ids) > 1
}
