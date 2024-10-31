package result

import "github.com/google/uuid"

type Result struct {
	ID          *uuid.UUID        `json:"id"`
	Title       string            `json:"title"`
	Signature   *string           `json:"signature"`
	TestOpsID   *int64            `json:"testops_id"`
	Execution   Execution         `json:"execution"`
	Fields      map[string]string `json:"fields"`
	Attachments []Attachment      `json:"attachments"`
	Steps       []Step            `json:"steps"`
	StepType    string            `json:"step_type,omitempty"`
	Params      map[string]string `json:"params"`
	ParamGroups [][]string        `json:"param_groups"`
	Relations   Relation          `json:"relations"`
	Muted       bool              `json:"muted"`
	Message     *string           `json:"message"`
	StartTime   *float64          `json:"start_time"`
	EndTime     *float64          `json:"end_time"`
	Duration    *float64          `json:"duration"`
}
