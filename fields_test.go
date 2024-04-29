package retraced

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected Fields
	}{
		{
			name:     "GQL format",
			data:     `[{"key": "abc", "value": "xyz"}]`,
			expected: Fields{"abc": "xyz"},
		},
		{
			name: "Go format",
			data: `{"abc": "xyz", "is_true": true}`,
			expected: Fields{
				"abc":     "xyz",
				"is_true": "true",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := Fields{}
			if err := json.Unmarshal([]byte(test.data), &fields); err != nil {
				t.Errorf("UnmarshalJSON failed: %v", err)
				return
			}
			assert.New(t).Equal(test.expected, fields, "Fields should be equal")
		})
	}
}
