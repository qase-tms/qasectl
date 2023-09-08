package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionTestPlanRunSummaries struct {
	Summaries []ActionTestPlanRunSummary
}

func (a *ActionTestPlanRunSummaries) TypeName() string {
	return "ActionTestPlanRunSummaries"
}

func (a *ActionTestPlanRunSummaries) Decode(m map[string]any) {
	a.Summaries = xcresult.DecodeArray[ActionTestPlanRunSummary](m["summaries"].(map[string]any))
}
