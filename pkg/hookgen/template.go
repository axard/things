package hookgen

import (
	"io"
	"strings"
	"text/template"
)

// Use only spaces for indentation
const hooktemplate = `// Code generated by hookgen; DO NOT EDIT.
// github.com/axard/things/cmd/hookgen

package {{.PackageName}}

import (
    {{- range .Imports }}
    "{{.}}"
    {{- end }}
    {{- with .Safe}}
    "sync"
    {{- end}}
)

type {{.HookName}} struct {
    list []*hooked
    {{- with .Safe}}
    m sync.Mutex
    {{- end}}
}

type hooked struct {
    {{.InterfaceName}}
}

type Cancel = func()

func (this *{{.HookName}}) Append(item {{.InterfaceName}}) Cancel {
    {{with .Safe -}}
    this.m.Lock()
    defer this.m.Unlock()

    {{end -}}
    hooked := &hooked{item}
    this.list = append(this.list, hooked)

    return func() { this.remove(hooked) }
}

func (this *{{.HookName}}) remove(hooked *hooked) {
    {{with .Safe -}}
    this.m.Lock()
    defer this.m.Unlock()

    {{end -}}
    for i := range this.list {
        if this.list[i] == hooked {
            this.list = append(this.list[:i], this.list[i+1:]...)
            break
        }
    }
}

func (this *{{.HookName}}) {{.MethodName}}({{join .MethodDeclArgs ", "}}) {
    {{with .Safe -}}
    this.m.Lock()
    defer this.m.Unlock()

    {{end -}}
    for _, hooked := range this.list {
        hooked.{{.MethodName}}({{join .MethodCallArgs ", "}})
    }
}
`

type HookTemplate struct {
	Imports        []string
	InterfaceName  string
	HookName       string
	PackageName    string
	MethodName     string
	MethodDeclArgs []string
	MethodCallArgs []string

	Safe bool
}

func (this HookTemplate) String() string {
	return hooktemplate
}

func (this HookTemplate) Write(w io.Writer) error {
	t, err := template.
		New("hooktemplate").
		Funcs(template.FuncMap{
			"join": strings.Join,
		}).
		Parse(hooktemplate)
	if err != nil {
		return err
	}

	return t.Execute(w, this)
}
