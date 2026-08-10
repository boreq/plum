package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMethod(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Valid bool
	}{
		{Name: "get", Value: "GET", Valid: true},
		{Name: "unknown_method", Value: "FROBNICATE", Valid: true},
		{Name: "empty", Value: "", Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			method, err := NewMethod(testCase.Value)
			if !testCase.Valid {
				require.EqualError(t, err, "method is empty")
				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.Value, method.String())
		})
	}
}
