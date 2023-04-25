package model

import (
	"github.com/qase-tms/qasectl/internal/xcresult"
	"github.com/qase-tms/qasectl/pkg"
)

type UniformTypeIdentifier string

const (
	PlainText UniformTypeIdentifier = "public.plain-text"
	PNG       UniformTypeIdentifier = "public.png"
)

type ActionTestAttachment struct {
	TypeIdentifier *UniformTypeIdentifier
	Name           *string
	PayloadRef     *Reference
}

func (a *ActionTestAttachment) TypeName() string {
	return "ActionTestAttachment"
}

func (a *ActionTestAttachment) Decode(m map[string]any) {
	if typeIdentifier, ok := m["typeIdentifier"].(map[string]any); ok {
		typeIdentifierName := typeIdentifier["name"].(map[string]any)
		a.TypeIdentifier = pkg.Ptr(UniformTypeIdentifier(xcresult.DecodeString(typeIdentifierName)))
	}

	name := m["name"].(map[string]any)
	a.Name = pkg.Ptr(xcresult.DecodeString(name))

	if payloadRef, ok := m["payloadRef"].(map[string]any); ok {
		a.PayloadRef = new(Reference)
		a.PayloadRef.Decode(payloadRef)
	}
}
