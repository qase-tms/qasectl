package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionTestPlanRunSummary struct {
	ActionAbstractTestSummary
	TestableSummaries []ActionTestableSummary
}

func (a *ActionTestPlanRunSummary) TypeName() string {
	return "ActionTestPlanRunSummary"
}

func (a *ActionTestPlanRunSummary) Decode(m map[string]any) {
	a.ActionAbstractTestSummary.Decode(m)
	a.TestableSummaries = xcresult.DecodeArray[ActionTestableSummary](m["testableSummaries"].(map[string]any))
}
