package libmockery

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Cleanup func() error

type OutputStreamProvider interface {
	GetWriter(iface *Interface) (io.Writer, error, Cleanup)
}

type StdoutStreamProvider struct {
}

func (ssp *StdoutStreamProvider) GetWriter(iface *Interface) (io.Writer, error, Cleanup) {
	return os.Stdout, nil, func() error { return nil }
}

type FileOutputStreamProvider struct {
	BaseDir string
}

func (fosp *FileOutputStreamProvider) GetWriter(iface *Interface) (io.Writer, error, Cleanup) {
	caseName := fosp.underscoreCaseName(iface.Name)
	path := filepath.Join(fosp.BaseDir, fosp.filename(caseName))
	os.MkdirAll(filepath.Dir(path), 0755)

	f, err := os.Create(path)
	if err != nil {
		return nil, err, nil
	}

	fmt.Printf("generating mock for interface: %s\n", iface.Name)
	return f, nil, func() error {
		return f.Close()
	}
}

func (fosp *FileOutputStreamProvider) filename(name string) string {
	return fmt.Sprintf("%s.go", name)
}

// shamelessly taken from http://stackoverflow.com/questions/1175208/elegant-python-function-to-convert-camelcase-to-camel-caseo
func (fosp *FileOutputStreamProvider) underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
