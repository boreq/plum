package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUri(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
	}{
		{Name: "root", Value: "/"},
		{Name: "query", Value: "/search?q=hello+world"},
		{Name: "not_normalized", Value: "/%2e%2e/etc/passwd"},
		{Name: "empty", Value: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			require.Equal(t, testCase.Value, NewUri(testCase.Value).String())
		})
	}
}
