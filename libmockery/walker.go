package libmockery

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Walker struct {
	BaseDir   string
	Interface string
}

type WalkerVisitor interface {
	VisitWalk(*Interface) error
}

// Walk returns true if a mock is generated.
func (w *Walker) Walk(visitor WalkerVisitor) bool {
	parser := NewParser()

	if err := w.doWalk(parser, w.BaseDir, visitor); err != nil {
		fmt.Printf("failed to walk directory %q: %v\n", w.BaseDir, err)
		os.Exit(1)
	}

	if err := parser.Load(); err != nil {
		fmt.Printf("parser had error while walking: %v\n", err)
		os.Exit(1)
	}

	filter := regexp.MustCompile(
		fmt.Sprintf(`^%s$`, w.Interface),
	)

	for _, iface := range parser.Interfaces() {
		if !filter.MatchString(iface.Name) {
			continue
		}

		if err := visitor.VisitWalk(iface); err != nil {
			fmt.Printf("visitor had error in walk for interface %q: %s\n", iface.Name, err)
			os.Exit(1)
		}

		return true
	}

	return false
}

// doWalk is a helper function that returns true as soon as the first mock is generated.
func (w *Walker) doWalk(p *Parser, dir string, visitor WalkerVisitor) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := file.Name()
		if strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "_") {
			continue
		}

		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, filename)

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		if err := p.Parse(path); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing file: ", err)
			return err
		}
	}

	return nil
}

type GeneratorVisitor struct {
	Comment           string
	OutputPackageName string
	OutputProvider    OutputStreamProvider
	ImportPrefix      string
}

func (gv *GeneratorVisitor) VisitWalk(iface *Interface) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to generated mock for '%s': %s\n", iface.Name, r)
			return
		}
	}()

	var out io.Writer

	out, err, closer := gv.OutputProvider.GetWriter(iface)
	if err != nil {
		fmt.Printf("Unable to get writer for %s: %s", iface.Name, err)
		os.Exit(1)
	}
	defer closer()

	gen := NewGenerator(iface, gv.OutputPackageName, gv.ImportPrefix)
	gen.GeneratePrologueComment(gv.Comment)
	gen.GeneratePrologue(gv.OutputPackageName)

	if err = gen.Generate(); err != nil {
		return err
	}

	if err = gen.Write(out); err != nil {
		return err
	}
	return nil
}
