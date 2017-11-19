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
	GetWriter(iface *Interface, pkg string) (io.Writer, error, Cleanup)
}

type StdoutStreamProvider struct {
}

func (ssp *StdoutStreamProvider) GetWriter(iface *Interface, pkg string) (io.Writer, error, Cleanup) {
	return os.Stdout, nil, func() error { return nil }
}

type FileOutputStreamProvider struct {
	BaseDir   string
	InPackage bool
	TestOnly  bool
	Case      string
}

func (fosp *FileOutputStreamProvider) GetWriter(iface *Interface, pkg string) (io.Writer, error, Cleanup) {
	var path string

	caseName := iface.Name
	if fosp.Case == "underscore" {
		caseName = fosp.underscoreCaseName(caseName)
	}

	if fosp.InPackage {
		path = filepath.Join(filepath.Dir(iface.Path), fosp.filename(caseName))
	} else {
		path = filepath.Join(fosp.BaseDir, fosp.filename(caseName))
		os.MkdirAll(filepath.Dir(path), 0755)
		pkg = filepath.Base(filepath.Dir(path))
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err, func() error { return nil }
	}

	fmt.Printf("Generating mock for: %s\n", iface.Name)
	return f, nil, func() error {
		return f.Close()
	}
}

func (fosp *FileOutputStreamProvider) filename(name string) string {
	if fosp.InPackage && fosp.TestOnly {
		return "mock_" + name + "_test.go"
	} else if fosp.InPackage {
		return "mock_" + name + ".go"
	} else if fosp.TestOnly {
		return name + "_test.go"
	}
	return name + ".go"
}

// shamelessly taken from http://stackoverflow.com/questions/1175208/elegant-python-function-to-convert-camelcase-to-camel-caseo
func (fosp *FileOutputStreamProvider) underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
