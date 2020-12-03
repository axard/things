package hookgen

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"path"

	"github.com/axard/things/pkg/formatter"
	"github.com/axard/things/pkg/resource"
)

type Hookgen struct {
	Src string
	Dst string

	Safe bool

	Formatter string
}

func (this *Hookgen) Generate(w io.Writer, i *types.Interface) error {
	ht := &HookTemplate{
		Imports:        this.methodImports(i),
		InterfaceName:  this.interfaceName(resource.Object(this.Src)),
		HookName:       this.hookName(resource.Object(this.Dst)),
		PackageName:    this.packageName(i.Method(0).Pkg().Name()),
		MethodName:     this.methodName(i),
		MethodDeclArgs: this.methodDeclArgs(i),
		MethodCallArgs: this.methodCallArgs(i),

		Safe: this.Safe,
	}

	buf := bytes.Buffer{}
	if err := ht.Write(&buf); err != nil {
		return err
	}

	formatter := this.formatter()

	formatted, err := formatter(buf.Bytes())
	if err != nil {
		return err
	}

	if _, err := w.Write(formatted); err != nil {
		return err
	}

	return nil
}

type formatterFunc func([]byte) ([]byte, error)

func (this *Hookgen) formatter() formatterFunc {
	switch this.Formatter {
	case "goimports":
		return formatter.Goimports
	case "noop":
		return formatter.Nofmt
	default:
		return formatter.Gofmt
	}
}

func (this *Hookgen) interfaceName(name string) string {
	if resource.Package(this.Src) == resource.Package(this.Dst) {
		return name
	}

	return path.Base(resource.Package(this.Src)) + "." + name
}

func (this *Hookgen) hookName(name string) string {
	if name == "" {
		return "Hook"
	}

	return name
}

func (this *Hookgen) packageName(name string) string {
	srcPkgName := path.Base(resource.Package(this.Src))
	dstPkgName := path.Base(resource.Package(this.Dst))

	if (srcPkgName == dstPkgName) && (name == "main") {
		return name
	}

	return dstPkgName
}

func (this *Hookgen) methodName(iface *types.Interface) string {
	return iface.Method(0).Name()
}

func (this *Hookgen) methodDeclArgs(iface *types.Interface) []string {
	args := []string{}

	sign := iface.Method(0).Type().(*types.Signature)
	prms := sign.Params()

	qualifier := func(p *types.Package) string {
		if this.Src != "" && this.Src == p.Path() {
			return ""
		}

		return p.Name()
	}

	for i := 0; i < prms.Len(); i++ {
		prm := prms.At(i)

		isLast := i == prms.Len()-1
		isSlice := prm.Type().String()[0:2] == "[]"
		isVariadic := sign.Variadic()

		n := prm.Name()
		if n == "" || n == "_" {
			n = fmt.Sprintf("arg%d", i)
		}

		t := types.TypeString(prm.Type(), qualifier)
		if isLast && isSlice && isVariadic {
			t = "..." + t[2:]
		}

		arg := n + " " + t

		args = append(args, arg)
	}

	return args
}

func (this *Hookgen) methodCallArgs(iface *types.Interface) []string {
	args := []string{}

	sign := iface.Method(0).Type().(*types.Signature)
	prms := sign.Params()

	for i := 0; i < prms.Len(); i++ {
		prm := prms.At(i)

		isLast := i == prms.Len()-1
		isSlice := prm.Type().String()[0:2] == "[]"
		isVariadic := sign.Variadic()

		n := prm.Name()
		if n == "" || n == "_" {
			n = fmt.Sprintf("arg%d", i)
		}

		arg := n
		if isLast && isSlice && isVariadic {
			arg = arg + "..."
		}

		args = append(args, arg)
	}

	return args
}

func (this *Hookgen) methodImports(iface *types.Interface) []string {
	importsMap := map[string]struct{}{}
	meth := iface.Method(0)
	sign := meth.Type().(*types.Signature)
	prms := sign.Params()

	if resource.Package(this.Src) != resource.Package(this.Dst) {
		importsMap[meth.Pkg().Path()] = struct{}{}
	}

	qualifier := func(p *types.Package) string {
		importsMap[p.Path()] = struct{}{}
		return ""
	}

	for i := 0; i < prms.Len(); i++ {
		types.TypeString(prms.At(i).Type(), qualifier)
	}

	imports := make([]string, 0, len(importsMap))
	for i := range importsMap {
		imports = append(imports, i)
	}

	return imports
}
