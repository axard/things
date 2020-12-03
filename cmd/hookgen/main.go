package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/axard/things/pkg/hookgen"
	"github.com/axard/things/pkg/resource"
	"golang.org/x/tools/go/packages"
)

type TFlags struct {
	Safe        bool
	ShowVersion bool

	Formatter string

	PathToSrc string
	PathToDst string

	File string
}

const (
	DirPermission  = 0755
	FilePermission = 0644
)

var (
	ErrEmptyPathToSrc   = errors.New("flag '-src' can't be empty")
	ErrEmptyPathToDst   = errors.New("flag '-dst' can't be empty")
	ErrInvalidPathToSrc = errors.New("flag '-src' can't be empty")
)

func (this *TFlags) Validate() error {
	if this.PathToSrc == "" {
		return ErrEmptyPathToSrc
	}

	if this.PathToDst == "" {
		return ErrEmptyPathToDst
	}

	if resource.Object(this.PathToSrc) == "" {
		return ErrInvalidPathToSrc
	}

	return nil
}

var (
	Version string = "unset"

	Flags TFlags
)

func init() {
	flag.BoolVar(&Flags.Safe, "safe", false, "use mutex to protect hook methods")
	flag.BoolVar(&Flags.ShowVersion, "version", false, "show the version for hoog")
	flag.StringVar(&Flags.Formatter, "fmt", "gofmt", "go pretty-printer: gofmt, goimports or noop (default gofmt)")
	flag.StringVar(&Flags.PathToSrc, "src", "", "path to interface like: /path/to/package.InterfaceName")
	flag.StringVar(&Flags.PathToDst, "dst", "", "path to hook like: /path/to/package[.HookName]")
	flag.StringVar(&Flags.File, "file", "generated.go", "name of generated file")

	flag.Parse()
}

func main() {
	if Flags.ShowVersion {
		fmt.Printf("hoog version: %s\n", Version)
		os.Exit(0)
	}

	if err := Flags.Validate(); err != nil {
		fmt.Printf("error: %s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	b := bytes.Buffer{}
	h := hookgen.Hookgen{
		Src:       Flags.PathToSrc,
		Dst:       Flags.PathToDst,
		Safe:      Flags.Safe,
		Formatter: Flags.Formatter,
	}

	iface, err := loadInterface(resource.Package(Flags.PathToSrc), resource.Object(Flags.PathToSrc))
	if err != nil {
		fatal(err)
	}

	if err := h.Generate(&b, iface); err != nil {
		fatal(err)
	}

	if err := os.MkdirAll(resource.Package(Flags.PathToDst), DirPermission); err != nil {
		fatal(err)
	}

	filename := filepath.Join(resource.Package(Flags.PathToDst), Flags.File)
	if err := ioutil.WriteFile(filename, b.Bytes(), FilePermission); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Printf("error: %s\n", err)
	os.Exit(1)
}

func loadInterface(dir, name string) (*types.Interface, error) {
	ps, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
		Dir:  dir,
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, fmt.Errorf("There are no packages in '%s'", dir)
	}

	if len(ps) > 1 {
		return nil, fmt.Errorf("Too many packages in '%s'", dir)
	}

	if errs := ps[0].Errors; len(errs) != 0 {
		errstr := make([]string, 0, len(errs))

		for _, err := range errs {
			errstr = append(errstr, err.Error())
		}

		return nil, fmt.Errorf("There are many errors: %s\n", strings.Join(errstr, "\n - "))
	}

	obj := ps[0].Types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("There is no interface with name '%s' in dir '%s'", name, dir)
	}

	if !types.IsInterface(obj.Type()) {
		return nil, fmt.Errorf("'%s' is not an interface; it is '%s'", name, obj.Type().String())
	}

	return obj.Type().Underlying().(*types.Interface), nil
}
