package libmockery

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type Parser interface {
	Parse(path string) error
	Load() error
	Find(name string) (*Interface, error)
	FindInPackage(name string, pkg *types.Package) *Interface
	Interfaces() []*Interface
}

type Parser2 struct {
	config  *packages.Config
	pkg     *packages.Package
	showLog bool
}

func NewParser2(verbose bool) *Parser2 {
	return &Parser2{
		config: &packages.Config{
			Mode:       packages.LoadTypes, // enough.
			Context:    nil,                // set a timeout?
			Dir:        "",                 // this dir?
			Env:        os.Environ(),       // same?
			BuildFlags: []string{},         // empty?
			Fset:       nil,                // defaults?
			ParseFile:  nil,                // defaults?
			Tests:      true,               // parse _test files?
		},
		pkg:     nil,     // gets set in Parse
		showLog: verbose, // show lots of trace logging
	}
}

func (p *Parser2) logf(format string, i ...interface{}) {
	if p.showLog {
		fmt.Printf("Parser2: "+format+"\n", i...)
	}
}

func stripMajor(module string) string {
	tokens := strings.Split(module, "/")
	last := tokens[len(tokens)-1]
	if verRe.MatchString(last) {
		return strings.Join(tokens[:len(tokens)-1], "/")
	}
	return strings.Join(tokens, "/")
}

func (p *Parser2) Parse(filename string) error {
	// filename is just a filename, in the current directory
	// from that filename, we determine the Go package that
	// the file exists in, relative to the nearest go.mod file
	// walking up the file tree.

	p.logf("Parse filename: %s", filename)

	pkgBase, pkgDir, err := pkgFromFile(filename)
	if err != nil {
		return err
	}

	p.logf("pkgBase: %s, pkgDir: %s\n", pkgBase, pkgDir)

	abs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	dir := filepath.Dir(abs)
	p.logf("abs dir of file: %s", dir)

	modFile, err := p.findModFile(dir)
	if err != nil {
		return err
	}

	p.logf("found mod: %s", modFile)

	module, err := moduleFromModFile(modFile)
	if err != nil {
		return err
	}

	p.logf("module: %s", module)
	moduleWithoutMajorVersion := stripMajor(module)
	p.logf("module no major: %s", moduleWithoutMajorVersion)

	i := strings.Index(abs, moduleWithoutMajorVersion)
	if i < 0 {
		return errors.New("module does not exist in abs")
	}

	rest := abs[1+i+len(module):]

	pkgs := filepath.Dir(rest)
	pkg := filepath.Join(module, pkgs)
	p.logf("determined package is: %s", pkg)

	// need to determine the package name in which filename exists
	// e.g. "github.com/a/b/c/d"
	loadedPkgs, err := packages.Load(p.config, pkg)
	if err != nil {
		return err
	}

	if len(loadedPkgs) == 0 {
		return errors.New("failed to load any package")
	}

	if len(loadedPkgs) > 1 {

		for _, loadedPkg := range loadedPkgs {
			if loadedPkg.ID == pkg {
				p.pkg = loadedPkg
				p.logf(
					"loaded package name: %s, id: %s",
					loadedPkg.Name,
					loadedPkg.ID,
				)
				return nil
			}
			p.logf(
				"ignoring package name: %s, id: %s",
				loadedPkg.Name,
				loadedPkg.ID,
			)
		}

		return errors.Errorf("no loaded packages match: %s", pkg)
	}

	p.pkg = loadedPkgs[0]
	p.logf("loaded only matching package name: %s, id: %s", p.pkg.Name, p.pkg.ID)
	return nil
}

var (
	pkgRe = regexp.MustCompile(`package[\s]+([\S]+)`)
	modRe = regexp.MustCompile(`module[\s]+([\S]+)`)
	verRe = regexp.MustCompile(`v[0-9]+`)
)

func moduleFromModFile(filename string) (string, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	content := string(bs)
	groups := modRe.FindStringSubmatch(content)
	if len(groups) != 2 {
		return "", errors.New("module declaration not found in go.mod file")
	}

	return groups[1], nil
}

func pkgFromFile(filename string) (string, string, error) {
	if filepath.Ext(filename) != ".go" {
		return "", "", errors.Errorf("%s is not a go file", filename)
	}

	// name of the package declared in filename
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", err
	}

	content := string(bs)
	groups := pkgRe.FindStringSubmatch(content)
	if len(groups) != 2 {
		return "", "", errors.New("package not found")
	}

	// name of the directory the base package is in
	abs, err := filepath.Abs(filename)
	if err != nil {
		return "", "", err
	}

	return groups[1], filepath.Base(abs), nil
}

func (p *Parser2) findModFile(path string) (string, error) {
	p.logf("findModFile in: %s", path)

	if path == "/" {
		return "", errors.New("go.mod did not exist in tree")
	}

	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, info := range infos {
		name := info.Name()
		if name == "go.mod" {
			return filepath.Join(path, name), nil
		}
	}

	upOne := filepath.Dir(path)
	p.logf("upOne: %s", upOne)

	return p.findModFile(upOne)
}

func (p *Parser2) Load() error {
	// loader version creates a map from abs filepath to declared interfaces
	// in that file
	p.logf("Load: go files: %s", p.pkg.GoFiles)

	return nil
}

func (p *Parser2) Find(path string) (*Interface, error) {
	panic("Find not implemented in v3")
}

func (p *Parser2) FindInPackage(name string, pkg *types.Package) *Interface {
	panic("FindInPackage not implemented in v3")
}

func (p *Parser2) Interfaces() []*Interface {
	p.logf("Interfaces")

	ifaces := make([]*Interface, 0, 10)
	scope := p.pkg.Types.Scope()
	for _, name := range scope.Names() {
		object := scope.Lookup(name)
		name := object.Name()
		typ := object.Type()
		id := object.Id()

		p.logf("Load: object name: %s, type: %s, id: %s", name, typ, id)

		isIface := types.IsInterface(typ)
		p.logf(" => is interface: %t", isIface)

		if isIface {
			ifaces = append(ifaces, &Interface{
				Name:      name,
				Type:      typ.Underlying().(*types.Interface),
				NamedType: typ.(*types.Named),
			})
		}
	}

	sort.Sort(sortableIFaceList(ifaces))
	return ifaces
}
