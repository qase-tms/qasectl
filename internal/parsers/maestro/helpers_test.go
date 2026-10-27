package maestro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type cmdSpec struct {
	Type       string
	Args       map[string]any
	Status     string
	Depth      int
	Ts, Dur    int64
	Err, Debug string
	Artifacts  []map[string]string
}

func header(name, file string, tags []string, props map[string]string, udid, shard string) []cmdSpec {
	config := map[string]any{"appId": "com.example.app", "name": name}
	if len(tags) > 0 {
		config["tags"] = tags
	}
	if len(props) > 0 {
		config["properties"] = props
	}
	return []cmdSpec{
		{
			Type: "defineVariablesCommand",
			Args: map[string]any{"env": map[string]any{
				"MAESTRO_FILENAME":    file,
				"MAESTRO_DEVICE_UDID": udid,
				"MAESTRO_SHARD_ID":    shard,
			}},
			Status: "COMPLETED", Ts: 1000,
		},
		{
			Type:   "applyConfigurationCommand",
			Args:   map[string]any{"config": config},
			Status: "COMPLETED", Ts: 1000,
		},
	}
}

func writeFlow(t *testing.T, root, dir string, cmds []cmdSpec) string {
	t.Helper()
	data := encodeCommands(t, cmds)
	flowDir := filepath.Join(root, dir)
	if err := os.MkdirAll(flowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(flowDir, "commands.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	return flowDir
}

func decodeCommands(t *testing.T, cmds []cmdSpec) []command {
	t.Helper()
	var out []command
	if err := json.Unmarshal(encodeCommands(t, cmds), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func encodeCommands(t *testing.T, cmds []cmdSpec) []byte {
	t.Helper()
	out := make([]map[string]any, 0, len(cmds))
	for i, c := range cmds {
		args := map[string]any{"optional": false}
		for k, v := range c.Args {
			args[k] = v
		}
		meta := map[string]any{
			"status":         c.Status,
			"timestamp":      c.Ts,
			"duration":       c.Dur,
			"sequenceNumber": i,
			"depth":          c.Depth,
		}
		if c.Err != "" || c.Debug != "" {
			meta["error"] = map[string]any{"message": c.Err, "debugMessage": c.Debug}
		}
		if len(c.Artifacts) > 0 {
			meta["artifacts"] = c.Artifacts
		}
		out = append(out, map[string]any{
			"command":  map[string]any{c.Type: args},
			"metadata": meta,
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeFile(t *testing.T, flowDir, rel string, size int) string {
	t.Helper()
	p := filepath.Join(flowDir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}
