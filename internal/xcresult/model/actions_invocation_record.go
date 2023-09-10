package model

import (
	"encoding/json"
	"github.com/qase-tms/qasectl/internal/xcresult"
)

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

func (a *ActionsInvocationRecord) UnmarshalJSON(bytes []byte) error {
	err := assertType(bytes, "ActionsInvocationRecord")
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, a)
}
