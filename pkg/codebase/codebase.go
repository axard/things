package codebase

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	Interface struct {
		Name   string
		Method *FunSignature
	}

	FunSignature struct {
		Name   string
		Params []*VarSignature
	}

	VarSignature struct {
		Name string
		Type string
	}

	// Should be here because of test
	dummy interface {
		Method(arg string)
	}
)

var (
	ErrEmptyName      = errors.New("empty name")
	ErrNotFound       = errors.New("not found in scope")
	ErrNotInterface   = errors.New("not interface")
	ErrEmptyInterface = errors.New("empty interface")
	ErrManyMethods    = errors.New("interface has more than one method")
	ErrHasResults     = errors.New("interface method has result")

	// Should be here because of test
	_ dummy = nil
)

func NewInterface(name string) *Interface {
	return &Interface{
		Name: name,
	}
}

func NewFunSignature(name string) *FunSignature {
	return &FunSignature{
		Name: name,
	}
}

func NewVarSignature(name string) *VarSignature {
	return &VarSignature{
		Name: name,
	}
}

func (this *Interface) LoadFrom(pkg *packages.Package) error {
	if this.Name == "" {
		return ErrEmptyName
	}

	record := pkg.Types.Scope().Lookup(this.Name)
	if record == nil {
		return ErrNotFound
	}

	if !types.IsInterface(record.Type()) {
		return ErrNotInterface
	}

	ifaceObj := record.Type().Underlying().(*types.Interface).Complete()
	if ifaceObj.NumMethods() == 0 {
		return ErrEmptyInterface
	}

	if ifaceObj.NumMethods() > 1 {
		return ErrManyMethods
	}

	funcObj := ifaceObj.Method(0)
	method := NewFunSignature(funcObj.Name())

	if err := method.LoadFrom(funcObj); err != nil {
		return err
	}

	this.Method = method

	return nil
}

func (this *FunSignature) LoadFrom(funcObj *types.Func) error {
	signObj := funcObj.Type().(*types.Signature)

	paramObjs := signObj.Params()

	for i := 0; i < paramObjs.Len(); i++ {
		paramObj := paramObjs.At(i)

		v := NewVarSignature(paramObj.Name())
		v.Type = paramObj.Type().String()

		isLastParam := i == paramObjs.Len()-1
		isArrayType := v.Type[0:2] == "[]"
		isVariadic := signObj.Variadic() && isLastParam && isArrayType

		if isVariadic {
			v.Type = "..." + v.Type[2:]
		}

		this.Params = append(this.Params, v)
	}

	resultObjs := signObj.Results()
	if resultObjs.Len() != 0 {
		return ErrHasResults
	}

	return nil
}
