package maestro

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	udidPro = "135E741D-8D9B-4E28-84C9-B4F31C81CFFA"
	udidStd = "607CC6EC-0E56-4457-B8AD-EA88129D2EEF"
)

const shardJUnit = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="Test Suite" device="iPhone 16 Pro - iOS 18.4 - ` + udidPro + `" tests="1">
    <testcase name="Probe" file="/abs/props.yaml"/>
  </testsuite>
  <testsuite name="Test Suite" device="iPhone 16 - iOS 18.4 - ` + udidStd + `" tests="1">
    <testcase name="Probe" file="/abs/props.yaml"/>
  </testsuite>
</testsuites>`

func writeXML(t *testing.T, root, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func shardFlow(name, udid, shard string) []cmdSpec {
	cmds := header(name, "f", nil, nil, udid, shard)
	return append(cmds, cmdSpec{Type: "tapOnElement", Status: "COMPLETED", Ts: 1100, Dur: 10})
}

func TestDeviceNames(t *testing.T) {
	root := t.TempDir()
	writeXML(t, root, "report.xml", shardJUnit)
	got := deviceNames(root)
	if len(got) != 2 || got[udidPro] != "iPhone 16 Pro - iOS 18.4" || got[udidStd] != "iPhone 16 - iOS 18.4" {
		t.Errorf("got %v", got)
	}
}

func TestDeviceNames_BareTestsuiteRoot(t *testing.T) {
	root := t.TempDir()
	writeXML(t, root, "sub.xml", `<testsuite device="Pixel 7 - Android 14 - emulator-5554"/>`)
	if got := deviceNames(root); got["emulator-5554"] != "Pixel 7 - Android 14" {
		t.Errorf("got %v", got)
	}
}

func TestDeviceNames_IgnoresForeignXML(t *testing.T) {
	root := t.TempDir()
	writeXML(t, root, "a.plist", `<?xml version="1.0"?><plist><dict/></plist>`)
	writeXML(t, root, "b.xml", `<?xml version="1.0"?><plist><dict/></plist>`)
	writeXML(t, root, "c.xml", `<testsuites><testsuite name="x"/></testsuites>`)
	writeXML(t, root, "d.xml", `not xml at all <<<`)
	if got := deviceNames(root); len(got) != 0 {
		t.Errorf("got %v", got)
	}
}

func TestParse_ShardAll(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "Probe-shard-1", shardFlow("Probe", udidPro, "1"))
	writeFlow(t, root, "Probe-shard-2", shardFlow("Probe", udidStd, "2"))
	writeXML(t, root, "report.xml", shardJUnit)

	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("got %d results", len(res))
	}
	want := []struct{ thread, device string }{
		{"shard-1", "iPhone 16 Pro - iOS 18.4"},
		{"shard-2", "iPhone 16 - iOS 18.4"},
	}
	for i, w := range want {
		if res[i].Execution.Thread == nil || *res[i].Execution.Thread != w.thread {
			t.Errorf("result %d thread %v", i, res[i].Execution.Thread)
		}
		if res[i].Params["device"] != w.device {
			t.Errorf("result %d params %v", i, res[i].Params)
		}
	}
	if *res[0].Signature == *res[1].Signature {
		t.Errorf("signatures equal: %s", *res[0].Signature)
	}
}

func TestParse_SingleShardNoThread(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "Probe", shardFlow("Probe", udidPro, "1"))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	if res[0].Execution.Thread != nil {
		t.Errorf("thread %q", *res[0].Execution.Thread)
	}
}

func TestParse_NoJUnitNoDeviceParam(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "Probe-shard-1", shardFlow("Probe", udidPro, "1"))
	writeFlow(t, root, "Probe-shard-2", shardFlow("Probe", udidStd, "2"))
	res, err := NewParser(root, "").Parse()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		if _, ok := r.Params["device"]; ok {
			t.Errorf("unexpected device param: %v", r.Params)
		}
	}
}

func TestParse_DuplicateRunsError(t *testing.T) {
	root := t.TempDir()
	writeFlow(t, root, "run1/Probe", shardFlow("Probe", udidPro, "1"))
	writeFlow(t, root, "run2/Probe", shardFlow("Probe", udidPro, "1"))
	_, err := NewParser(root, "").Parse()
	if err == nil || !strings.Contains(err.Error(), "several Maestro runs") {
		t.Errorf("err = %v", err)
	}
}
