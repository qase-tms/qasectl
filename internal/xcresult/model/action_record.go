package model

import "encoding/json"

type ActionRecord struct {
	result interface{}
}

func (a *ActionRecord) UnmarshalJSON(bytes []byte) error {
	err := assertType(bytes, "ActionRecord")
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, a)
}
