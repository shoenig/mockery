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
	Recursive bool
	Filter    *regexp.Regexp
	LimitOne  bool
}

type WalkerVisitor interface {
	VisitWalk(*Interface) error
}

func (w *Walker) Walk(visitor WalkerVisitor) (generated bool) {
	parser := NewParser()
	w.doWalk(parser, w.BaseDir, visitor)

	err := parser.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking: %v\n", err)
		os.Exit(1)
	}

	for _, iface := range parser.Interfaces() {
		if !w.Filter.MatchString(iface.Name) {
			continue
		}
		err := visitor.VisitWalk(iface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %s\n", iface.Name, err)
			os.Exit(1)
		}
		generated = true
		if w.LimitOne {
			return
		}
	}

	return
}

func (w *Walker) doWalk(p *Parser, dir string, visitor WalkerVisitor) (generated bool) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			if w.Recursive {
				generated = w.doWalk(p, path, visitor) || generated
				if generated && w.LimitOne {
					return
				}
			}
			continue
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		err = p.Parse(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing file: ", err)
			return
		}
	}

	return
}

type GeneratorVisitor struct {
	InPackage bool
	Note      string
	Osp       OutputStreamProvider
	// The name of the output package, if InPackage is false (defaults to "mocks")
	PackageName string
}

func (gv *GeneratorVisitor) VisitWalk(iface *Interface) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to generated mock for '%s': %s\n", iface.Name, r)
			return
		}
	}()

	var out io.Writer
	var pkg string

	if gv.InPackage {
		pkg = iface.Path
	} else {
		pkg = gv.PackageName
	}

	out, err, closer := gv.Osp.GetWriter(iface, pkg)
	if err != nil {
		fmt.Printf("Unable to get writer for %s: %s", iface.Name, err)
		os.Exit(1)
	}
	defer closer()

	gen := NewGenerator(iface, pkg, gv.InPackage)
	gen.GeneratePrologueNote(gv.Note)
	gen.GeneratePrologue(pkg)

	err = gen.Generate()
	if err != nil {
		return err
	}

	err = gen.Write(out)
	if err != nil {
		return err
	}
	return nil
}
