package model

import "github.com/qase-tms/qasectl/internal/xcresult"

type ActionTestSummaryIdentifiableObject interface {
	xcresult.Decoder
	ID() *string
}
