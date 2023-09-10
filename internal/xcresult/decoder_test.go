package xcresult

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_types(t *testing.T) {
	t.Run("just a _type", func(t *testing.T) {
		m := map[string]interface{}{
			"_type": map[string]interface{}{
				"_name": "String",
			},
			"_value": "Some Value",
		}

		assert.Equal(t, []string{"String"}, types(m))
	})
	t.Run("_supertype presented", func(t *testing.T) {
		m := map[string]interface{}{
			"_type": map[string]interface{}{
				"_name": "ActionTestSummaryGroup",
				"_supertype": map[string]interface{}{
					"_name": "ActionTestSummaryIdentifiableObject",
					"_supertype": map[string]interface{}{
						"_name": "ActionAbstractTestSummary",
					},
				},
			},
		}
		expected := []string{
			"ActionTestSummaryGroup",
			"ActionTestSummaryIdentifiableObject",
			"ActionAbstractTestSummary",
		}

		assert.Equal(t, expected, types(m))
	})
}
