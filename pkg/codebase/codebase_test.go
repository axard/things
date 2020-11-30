package codebase

import (
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestInterface_LoadFrom(t *testing.T) {
	ps, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportsFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedModule,
	})

	if err != nil {
		t.Errorf("Error on package load: %s", err)
	}

	if len(ps) == 0 {
		t.Errorf("No packages found")
	}

	if len(ps) > 1 {
		t.Errorf("More than one package found")
	}

	this := &Interface{
		Name: "dummy",
	}

	if err := this.LoadFrom(ps[0]); err != nil {
		t.Errorf("Interface.LoadFrom() error = %v", err)
	}

	if len(this.Method.Params) != 1 {
		t.Errorf("Invalid method params number: got = %v; want = %v", len(this.Method.Params), 1)
	}

	if this.Method.Params[0].Name != "arg" {
		t.Errorf("Invalid param name: got = %s; want = %s", this.Method.Params[0].Name, "arg")
	}

	if this.Method.Params[0].Type != "string" {
		t.Errorf("Invalid param type: got = %s; want = %s", this.Method.Params[0].Type, "string")
	}
}
