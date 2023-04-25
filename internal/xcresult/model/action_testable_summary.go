package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
)

type ActionTestableSummary struct {
	ActionAbstractTestSummary
	Tests []xcresult.Decoder
}

func (a *ActionTestableSummary) TypeName() string {
	return "ActionTestableSummary"
}

func (a *ActionTestableSummary) Decode(m map[string]any) {
	a.ActionAbstractTestSummary.Decode(m)

	tests := m["tests"].(map[string]any)
	a.Tests = xcresult.DecodeVarArray(tests, func(typ string) xcresult.Decoder {
		switch typ {
		case "ActionTestSummaryGroup":
			return new(ActionTestSummaryGroup)
		}

		return nil
	})
}
