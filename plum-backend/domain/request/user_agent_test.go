package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUserAgent(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
	}{
		{Name: "browser", Value: "Mozilla/5.0 (X11; Linux x86_64) Firefox/143.0"},
		{Name: "dash", Value: "-"},
		{Name: "empty", Value: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			require.Equal(t, testCase.Value, NewUserAgent(testCase.Value).String())
		})
	}
}
