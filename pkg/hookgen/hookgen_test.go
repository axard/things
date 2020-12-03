package hookgen

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

func mustLoadInterface(dir, name string) *types.Interface {
	ps, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
		Dir:  dir,
	})
	if err != nil {
		panic(err)
	}

	if len(ps) == 0 {
		panic(fmt.Sprintf("There are no packages in '%s'", dir))
	}

	if len(ps) > 1 {
		panic(fmt.Sprintf("Too many packages in '%s'", dir))
	}

	if errs := ps[0].Errors; len(errs) != 0 {
		errstr := make([]string, 0, len(errs)+1)
		errstr = append(errstr, "there are many errors:")

		for _, err := range errs {
			errstr = append(errstr, err.Error())
		}

		panic(strings.Join(errstr, "\n - "))
	}

	obj := ps[0].Types.Scope().Lookup(name)
	if obj == nil {
		panic(fmt.Sprintf("There is no interface with name '%s' in dir '%s'", name, dir))
	}

	if !types.IsInterface(obj.Type()) {
		panic(fmt.Sprintf("'%s' is not an interface; it is '%s'", name, obj.Type().String()))
	}

	return obj.Type().Underlying().(*types.Interface)
}

func TestHookgen_Generate(t *testing.T) {
	type fields struct {
		SrcPkg    string
		DstPkg    string
		Safe      bool
		Formatter string
	}
	type args struct {
		i *types.Interface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Hookgen{
				Src:       tt.fields.SrcPkg,
				Dst:       tt.fields.DstPkg,
				Safe:      tt.fields.Safe,
				Formatter: tt.fields.Formatter,
			}
			w := &bytes.Buffer{}
			if err := this.Generate(w, tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Hookgen.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Hookgen.Generate() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestHookgen_methodName(t *testing.T) {
	type fields struct {
		SrcPkg    string
		DstPkg    string
		Safe      bool
		Formatter string
	}
	type args struct {
		iface *types.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface0"),
			},
			want: "Method",
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface1"),
			},
			want: "Method",
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface2"),
			},
			want: "Method",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Hookgen{
				Src:       tt.fields.SrcPkg,
				Dst:       tt.fields.DstPkg,
				Safe:      tt.fields.Safe,
				Formatter: tt.fields.Formatter,
			}
			if got := this.methodName(tt.args.iface); got != tt.want {
				t.Errorf("Hookgen.methodName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHookgen_methodDeclArgs(t *testing.T) {
	type fields struct {
		SrcPkg    string
		DstPkg    string
		Safe      bool
		Formatter string
	}
	type args struct {
		iface *types.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface0"),
			},
			want: []string{"arg0 interface{}"},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface1"),
			},
			want: []string{"arg0 ...interface{}"},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface2"),
			},
			want: []string{"s string", "arg1 interface{}"},
		},
		{
			name: "",
			fields: fields{
				SrcPkg: "github.com/axard/things/pkg/hookgen/internal/instance",
			},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface3"),
			},
			want: []string{"i int", "s Struct"},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface3"),
			},
			want: []string{"i int", "s instance.Struct"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Hookgen{
				Src:       tt.fields.SrcPkg,
				Dst:       tt.fields.DstPkg,
				Safe:      tt.fields.Safe,
				Formatter: tt.fields.Formatter,
			}
			if got := this.methodDeclArgs(tt.args.iface); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hookgen.methodDeclArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHookgen_methodCallArgs(t *testing.T) {
	type fields struct {
		SrcPkg    string
		DstPkg    string
		Safe      bool
		Formatter string
	}
	type args struct {
		iface *types.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface0"),
			},
			want: []string{"arg0"},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface1"),
			},
			want: []string{"arg0..."},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface2"),
			},
			want: []string{"s", "arg1"},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface3"),
			},
			want: []string{"i", "s"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Hookgen{
				Src:       tt.fields.SrcPkg,
				Dst:       tt.fields.DstPkg,
				Safe:      tt.fields.Safe,
				Formatter: tt.fields.Formatter,
			}
			if got := this.methodCallArgs(tt.args.iface); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hookgen.methodCallArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHookgen_methodImports(t *testing.T) {
	type fields struct {
		SrcPkg    string
		DstPkg    string
		Safe      bool
		Formatter string
	}
	type args struct {
		iface *types.Interface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface0"),
			},
			want: []string{},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface1"),
			},
			want: []string{},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface2"),
			},
			want: []string{},
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				iface: mustLoadInterface("./internal/instance", "Interface3"),
			},
			want: []string{"github.com/axard/things/pkg/hookgen/internal/instance"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Hookgen{
				Src:       tt.fields.SrcPkg,
				Dst:       tt.fields.DstPkg,
				Safe:      tt.fields.Safe,
				Formatter: tt.fields.Formatter,
			}
			if got := this.methodImports(tt.args.iface); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hookgen.methodImports() = %v, want %v", got, tt.want)
			}
		})
	}
}
