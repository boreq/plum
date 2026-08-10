package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVersion(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Valid bool
	}{
		{Name: "http_1_1", Value: "HTTP/1.1", Valid: true},
		{Name: "http_2", Value: "HTTP/2.0", Valid: true},
		{Name: "empty", Value: "", Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			version, err := NewVersion(testCase.Value)
			if !testCase.Valid {
				require.EqualError(t, err, "version is empty")
				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.Value, version.String())
		})
	}
}
