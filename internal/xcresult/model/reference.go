package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
	"github.com/qase-tms/qasectl/pkg"
)

type Reference struct {
	ObjID      string
	TargetType *string
}

func (r *Reference) TypeName() string {
	return "Reference"
}

func (r *Reference) Decode(m map[string]any) {
	if targetType, ok := m["targetType"].(map[string]any); ok {
		targetTypeName := targetType["name"].(map[string]any)
		r.TargetType = pkg.Ptr(xcresult.DecodeString(targetTypeName))
	}

	r.ObjID = xcresult.DecodeString(m["id"].(map[string]any))
}
