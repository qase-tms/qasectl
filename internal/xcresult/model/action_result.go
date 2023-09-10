package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionResult struct {
	TestsRef *Reference
}

func (a *ActionResult) TypeName() string {
	return "ActionResult"
}

func (a *ActionResult) Decode(m map[string]any) {
	v, ok := m["testsRef"].(map[string]any)
	if !ok {
		return
	}

	testsRef := xcresult.DecodeObject[Reference](v)
	a.TestsRef = &testsRef
}
