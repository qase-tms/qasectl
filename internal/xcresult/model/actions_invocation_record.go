package model

import "encoding/json"

type ActionsInvocationRecord struct {
	Actions []ActionRecord `json:"actions"`
}

func (a *ActionsInvocationRecord) UnmarshalJSON(bytes []byte) error {
	err := assertType(bytes, "ActionsInvocationRecord")
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, a)
}
