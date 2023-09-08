package model

import (
	"fmt"
	"github.com/qase-tms/qasectl/internal/xcresult"
)

func factory(typ string) xcresult.Decoder {
	switch typ {
	case "ActionTestSummaryGroup":
		return new(ActionTestSummaryGroup)
	case "ActionTestMetadata":
		return new(ActionTestMetadata)
	default:
		panic(fmt.Errorf("factory unknown type %q", typ))
	}
}
