package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStatus(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Valid bool
	}{
		{Name: "ok", Value: "200", Valid: true},
		{Name: "not_a_number", Value: "-", Valid: true},
		{Name: "empty", Value: "", Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			status, err := NewStatus(testCase.Value)
			if !testCase.Valid {
				require.EqualError(t, err, "status is empty")
				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.Value, status.String())
		})
	}
}
