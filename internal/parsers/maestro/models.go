package maestro

import (
	"encoding/json"
	"fmt"
)

type command struct {
	Type      string
	Args      json.RawMessage
	Label     string
	Optional  bool
	Status    string
	Timestamp int64
	Duration  int64
	Seq       int
	Depth     int
	Error     *commandError
	Artifacts []artifact
}

type commandError struct {
	Message      string `json:"message"`
	DebugMessage string `json:"debugMessage"`
}

type artifact struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type flow struct {
	Dir        string
	Name       string
	FileName   string
	UDID       string
	ShardID    string
	Tags       []string
	Properties map[string]string
	Commands   []command
}

type rawCommand struct {
	Command  map[string]json.RawMessage `json:"command"`
	Metadata struct {
		Status         string        `json:"status"`
		Timestamp      int64         `json:"timestamp"`
		Duration       int64         `json:"duration"`
		SequenceNumber int           `json:"sequenceNumber"`
		Depth          int           `json:"depth"`
		Error          *commandError `json:"error"`
		Artifacts      []artifact    `json:"artifacts"`
	} `json:"metadata"`
}

func (c *command) UnmarshalJSON(data []byte) error {
	var raw rawCommand
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw.Command) != 1 {
		return fmt.Errorf("command object must have exactly one key, got %d", len(raw.Command))
	}
	for k, v := range raw.Command {
		c.Type = k
		c.Args = v
	}

	var common struct {
		Label    string `json:"label"`
		Optional bool   `json:"optional"`
	}
	// Args of some commands are not objects; label/optional are best-effort.
	_ = json.Unmarshal(c.Args, &common)
	c.Label = common.Label
	c.Optional = common.Optional

	m := raw.Metadata
	c.Status = m.Status
	c.Timestamp = m.Timestamp
	c.Duration = m.Duration
	c.Seq = m.SequenceNumber
	c.Depth = m.Depth
	c.Error = m.Error
	c.Artifacts = m.Artifacts
	return nil
}
