package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionsInvocationRecord struct {
	Actions []ActionRecord `json:"actions"`
}

func (a *ActionsInvocationRecord) TypeName() string {
	return "ActionsInvocationRecord"
}

func (a *ActionsInvocationRecord) Decode(m map[string]any) {
	actions := m["actions"].(map[string]any)
	a.Actions = xcresult.DecodeArray[ActionRecord](actions)
}
