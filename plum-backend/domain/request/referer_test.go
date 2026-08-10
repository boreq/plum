package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReferer(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
	}{
		{Name: "url", Value: "https://example.com/"},
		{Name: "dash", Value: "-"},
		{Name: "empty", Value: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			require.Equal(t, testCase.Value, NewReferer(testCase.Value).String())
		})
	}
}
