package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessCodeBaseImageStreamNameConvention(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input string
		want  string
	}{
		"one dot":    {input: "feature.one", want: "feature-one"},
		"one slash":  {input: "feature/one", want: "feature-one"},
		"mixed":      {input: "feature/one.two", want: "feature-one-two"},
		"mixed long": {input: "qa/feature//one.two.", want: "qa-feature--one-two-"},
		"no changes": {input: "feature-one", want: "feature-one"},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := ProcessCodeBaseImageStreamNameConvention(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
