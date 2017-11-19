package libmockery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilenameBare(t *testing.T) {
	out := FileOutputStreamProvider{}
	require.Equal(t, "name.go", out.filename("name"))
}

func TestUnderscoreCaseName(t *testing.T) {
	require.Equal(t, "notify_event", (&FileOutputStreamProvider{}).underscoreCaseName("NotifyEvent"))
	require.Equal(t, "repository", (&FileOutputStreamProvider{}).underscoreCaseName("Repository"))
	require.Equal(t, "http_server", (&FileOutputStreamProvider{}).underscoreCaseName("HTTPServer"))
	require.Equal(t, "awesome_http_server", (&FileOutputStreamProvider{}).underscoreCaseName("AwesomeHTTPServer"))
	require.Equal(t, "csv", (&FileOutputStreamProvider{}).underscoreCaseName("CSV"))
	require.Equal(t, "position0_size", (&FileOutputStreamProvider{}).underscoreCaseName("Position0Size"))
}
