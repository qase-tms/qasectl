package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
	"github.com/qase-tms/qasectl/pkg"
)

type ActionAbstractTestSummary struct {
	Identifier *string
	Name       *string
}

func (a *ActionAbstractTestSummary) ID() *string {
	return a.Identifier
}

func (a *ActionAbstractTestSummary) TypeName() string {
	// but actually it's not used anywhere
	return "ActionAbstractTestSummary"
}

func (a *ActionAbstractTestSummary) Decode(m map[string]any) {
	name := m["name"].(map[string]any)
	a.Name = pkg.Ptr(xcresult.DecodeString(name))

	if identifier, ok := m["identifier"].(map[string]any); ok {
		a.Identifier = pkg.Ptr(xcresult.DecodeString(identifier))
	}
}
