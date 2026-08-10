package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRemoteAddress(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Valid bool
	}{
		{Name: "ipv4", Value: "1.2.3.4", Valid: true},
		{Name: "ipv6", Value: "2001:db8::1", Valid: true},
		{Name: "dash", Value: "-", Valid: true},
		{Name: "empty", Value: "", Valid: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			remoteAddress, err := NewRemoteAddress(testCase.Value)
			if !testCase.Valid {
				require.EqualError(t, err, "remote address is empty")
				return
			}

			require.NoError(t, err)
			require.Equal(t, testCase.Value, remoteAddress.String())
		})
	}
}
