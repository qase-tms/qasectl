package model

type ActionTestMetadata struct {
	ActionAbstractTestSummary
	SummaryRef *Reference
}

func (a *ActionTestMetadata) TypeName() string {
	return "ActionTestMetadata"
}

func (a *ActionTestMetadata) Decode(m map[string]any) {
	a.ActionAbstractTestSummary.Decode(m)

	if summaryRef, ok := m["summaryRef"].(map[string]any); ok {
		a.SummaryRef = new(Reference)
		a.SummaryRef.Decode(summaryRef)
	}
}
