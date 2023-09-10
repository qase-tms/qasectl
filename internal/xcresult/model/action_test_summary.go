package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
)

type TestStatus string

const (
	Failure  TestStatus = "Failure"
	Success  TestStatus = "Success"
	Undefine TestStatus = "Undefine"
)

type ActionTestSummary struct {
	ActionAbstractTestSummary
	Duration          float64
	Status            TestStatus
	ActivitySummaries []ActionTestActivitySummary
}

func (a *ActionTestSummary) TypeName() string {
	return "ActionTestSummary"
}

func (a *ActionTestSummary) Decode(m map[string]any) {
	a.ActionAbstractTestSummary.Decode(m)

	if duration, ok := m["duration"].(map[string]any); ok {
		a.Duration = xcresult.DecodeDouble(duration)
	}

	if status, ok := m["testStatus"].(map[string]any); ok {
		a.Status = TestStatus(xcresult.DecodeString(status))
	}

	if activitySummaries, ok := m["activitySummaries"].(map[string]any); ok {
		a.ActivitySummaries = xcresult.DecodeArray[ActionTestActivitySummary](activitySummaries)
	}
}
