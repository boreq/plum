package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBodyBytesSent(t *testing.T) {
	testCases := []struct {
		Name  string
		Value int
		Valid bool
	}{
		{Name: "positive", Value: 123, Valid: true},
		{Name: "zero", Value: 0, Valid: true},
		{Name: "negative", Value: -1, Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			bodyBytesSent, err := NewBodyBytesSent(testCase.Value)
			if !testCase.Valid {
				require.EqualError(t, err, "body bytes sent is negative")
				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.Value, bodyBytesSent.Int())
		})
	}
}
