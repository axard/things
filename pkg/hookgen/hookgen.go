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
		PackageName:    path.Base(resource.Package(this.Dst)),
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
	if this.Formatter == "goimports" {
		return formatter.Goimports
	}

	return formatter.Gofmt
}

func (this *Hookgen) interfaceName(name string) string {
	if this.Src == this.Dst {
		return name
	}

	return path.Base(this.Src) + "." + name
}

func (this *Hookgen) hookName(name string) string {
	if name == "" {
		return "Hook"
	}

	return name
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
	imports := []string{}

	qualifier := func(p *types.Package) string {
		imports = append(imports, p.Path())
		return ""
	}

	sign := iface.Method(0).Type().(*types.Signature)
	prms := sign.Params()

	for i := 0; i < prms.Len(); i++ {
		types.TypeString(prms.At(i).Type(), qualifier)
	}

	return imports
}
