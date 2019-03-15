package libmockery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testFile string
var testFile2 string

func init() {
	testFile = getFixturePath("requester.go")
	testFile2 = getFixturePath("requester2.go")
}

func TestFileParse(t *testing.T) {
	parser := NewParser2(true)

	err := parser.Parse(testFile)
	require.NoError(t, err)

	err = parser.Load()
	require.NoError(t, err)

	// method not implemented in v3
	// node, err := parser.Find("Requester")
	// require.NoError(t, err)
	// require.NotNil(t, node)
}

func testParse(t *testing.T, parser Parser, a, b, c string) {
	err := parser.Parse(getFixturePath(a, b, c))
	require.NoError(t, err)
}

func TestBuildTagInFilename(t *testing.T) {
	parser := NewParser2(true)

	// Include the major OS values found on https://golang.org/dl/ so we're likely to match
	// anywhere the test is executed.
	testParse(t, parser, "buildtag", "filename", "iface_windows.go")
	testParse(t, parser, "buildtag", "filename", "iface_linux.go")
	testParse(t, parser, "buildtag", "filename", "iface_darwin.go")
	testParse(t, parser, "buildtag", "filename", "iface_freebsd.go")

	err := parser.Load()
	require.NoError(t, err) // Expect "redeclared in this block" if tags aren't respected

	nodes := parser.Interfaces()
	require.Equal(t, 1, len(nodes))
	require.Equal(t, "IfaceWithBuildTagInFilename", nodes[0].Name)
}

func TestBuildTagInComment(t *testing.T) {
	parser := NewParser2(true)

	// Include the major OS values found on https://golang.org/dl/ so we're likely to match
	// anywhere the test is executed.
	testParse(t, parser, "buildtag", "comment", "windows_iface.go")
	testParse(t, parser, "buildtag", "comment", "linux_iface.go")
	testParse(t, parser, "buildtag", "comment", "darwin_iface.go")
	testParse(t, parser, "buildtag", "comment", "freebsd_iface.go")

	err := parser.Load()
	require.NoError(t, err) // Expect "redeclared in this block" if tags aren't respected

	nodes := parser.Interfaces()
	require.Equal(t, 1, len(nodes))
	require.Equal(t, "IfaceWithBuildTagInComment", nodes[0].Name)
}
