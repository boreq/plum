package request

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	remoteAddress, err := NewRemoteAddress("1.2.3.4")
	require.NoError(t, err)

	timestamp := time.Date(2026, time.August, 10, 12, 0, 0, 0, time.UTC)

	method, err := NewMethod("GET")
	require.NoError(t, err)

	uri := NewUri("/index.html")

	version, err := NewVersion("HTTP/1.1")
	require.NoError(t, err)

	status, err := NewStatus("200")
	require.NoError(t, err)

	bodyBytesSent, err := NewBodyBytesSent(123)
	require.NoError(t, err)

	referer := NewReferer("https://example.com/")
	userAgent := NewUserAgent("Mozilla/5.0 (X11; Linux x86_64) Firefox/143.0")

	request := NewRequest(
		remoteAddress,
		timestamp,
		method,
		uri,
		version,
		status,
		bodyBytesSent,
		referer,
		userAgent,
	)

	require.Equal(t, remoteAddress, request.RemoteAddress())
	require.Equal(t, timestamp, request.Timestamp())
	require.Equal(t, method, request.Method())
	require.Equal(t, uri, request.Uri())
	require.Equal(t, version, request.Version())
	require.Equal(t, status, request.Status())
	require.Equal(t, bodyBytesSent, request.BodyBytesSent())
	require.Equal(t, referer, request.Referer())
	require.Equal(t, userAgent, request.UserAgent())
}
