package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionTestSummaryGroup struct {
	ActionAbstractTestSummary
	Subtests []xcresult.Decoder
}

func (a *ActionTestSummaryGroup) TypeName() string {
	return "ActionTestSummaryGroup"
}

func (a *ActionTestSummaryGroup) Decode(m map[string]any) {
	a.ActionAbstractTestSummary.Decode(m)

	if subtests, ok := m["subtests"].(map[string]any); ok {
		a.Subtests = xcresult.DecodeVarArray(subtests, factory)
	}
}
