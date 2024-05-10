package _import

import "time"

type Result struct {
	ID          *int64
	Title       string
	Signature   *string
	TestOpsID   *int64
	Execution   Execution
	Fields      map[string]string
	Attachments []Attachment
	Steps       []Step
	StepType    string
	Params      map[string]string
	Relations   Relation
	Muted       bool
	Message     *string
	StartTime   *time.Time
	EndTime     *time.Time
	Duration    *time.Duration
}
