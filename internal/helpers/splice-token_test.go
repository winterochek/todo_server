package h

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpliceToken(t *testing.T) {
	tests := []struct {
		header   string
		expected string
		err      error
	}{
		// {
		// 	header:   "",
		// 	expected: "",
		// 	err:      ErrNoAuthHeader,
		// },
		// {
		// 	header:   "Bearer ",
		// 	expected: "",
		// 	err:      ErrWrongAuthHeaderType,
		// },
		{
			header:   "Bearer your-valid-token",
			expected: "your-valid-token",
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.header, func(t *testing.T) {
			token, err := SpliceToken(test.header)

			if err != nil {
				assert.ErrorIs(t, err, test.err)
			}

			if err == nil {
				assert.Equal(t, token, test.expected)
			}

		})

	}
}
