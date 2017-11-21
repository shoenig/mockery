package libmockery

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	file1 = `
package foo

type Bar struct { }
`

	file2 = `
package foo

type Baz interface {}
`
)

func makeFiles(t *testing.T) string {
	root, err := ioutil.TempDir("", "hasher-")
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(root, "file1.go"), []byte(file1), 0700)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(root, "file2.go"), []byte(file2), 0700)
	require.NoError(t, err)

	return root
}

func Test_Hash(t *testing.T) {
	root := makeFiles(t)

	h := NewHasher(root)

	content, err := h.Hash()
	require.NoError(t, err)

	require.Equal(t, 2, len(content))
}

func Test_Same_ok(t *testing.T) {
	a := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "bbbbbbbbb",
	}

	b := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "bbbbbbbbb",
	}

	err := Same(a, b)
	require.NoError(t, err)
}

func Test_Same_numFiles(t *testing.T) {
	a := Content{
		"/foo/bar":  "aaaaaaaaa",
		"/foo/baz":  "bbbbbbbbb",
		"/foo/derp": "cccccccc",
	}

	b := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "bbbbbbbbb",
	}

	err := Same(a, b)
	require.Error(t, err)
}

func Test_Same_differentFIles(t *testing.T) {
	a := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "bbbbbbbbb",
	}

	b := Content{
		"/foo/something": "aaaaaaaaa",
		"/foo/baz":       "bbbbbbbbb",
	}

	err := Same(a, b)
	require.Error(t, err)
}

func Test_Same_content(t *testing.T) {
	a := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "xxxxxxxxx",
	}

	b := Content{
		"/foo/bar": "aaaaaaaaa",
		"/foo/baz": "bbbbbbbbb",
	}

	err := Same(a, b)
	require.Error(t, err)
}
