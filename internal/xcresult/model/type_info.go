package model

import (
	"encoding/json"
	"fmt"
)

type typeInfo struct {
	Type struct {
		Name string `json:"_name"`
	} `json:"_type"`
}

func assertType(bytes []byte, expectedType string) error {
	var info typeInfo

	err := json.Unmarshal(bytes, &info)
	if err != nil {
		return err
	}

	if info.Type.Name != expectedType {
		return fmt.Errorf("underlying type is %q (expected %q)", info.Type.Name, expectedType)
	}

	return nil
}
