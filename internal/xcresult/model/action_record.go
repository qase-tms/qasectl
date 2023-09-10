package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionRecord struct {
	Result ActionResult
}

func (a *ActionRecord) TypeName() string {
	return "ActionRecord"
}

func (a *ActionRecord) Decode(m map[string]any) {
	a.Result = xcresult.DecodeObject[ActionResult](m["actionResult"].(map[string]any))
}

func (a *ActionRecord) UnmarshalJSON(bytes []byte) error {
	err := assertType(bytes, "ActionRecord")
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, a)
}
