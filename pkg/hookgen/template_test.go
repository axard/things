package hookgen

import (
	"bytes"
	"testing"
)

func TestHookTemplate_Execute(t *testing.T) {
	tests := []struct {
		name    string
		this    HookTemplate
		wantW   string
		wantErr bool
	}{
		{
			name: "Template{ Safe: false } doesn't use sync.Mutex",
			this: HookTemplate{
				Imports: []string{
					"io",
				},
				InterfaceName: "Callback",
				PackageName:   "cbhook",
				MethodName:    "Call",
				MethodDeclArgs: []string{
					"arg1 io.Writer",
					"arg2 ...interface{}",
				},
				MethodCallArgs: []string{
					"arg1",
					"arg2...",
				},
				Safe: false,
			},
			// Use only spaces for indentation
			wantW: `// Code generated by hookgen; DO NOT EDIT.
// github.com/axard/things/cmd/hookgen

package cbhook

import (
    "io"
)

type Hook struct {
    list []*hookedItem
}

type hooked struct {
    Callback
}

type Cancel = func()

func (this *Hook) Append(item Callback) Cancel {
    hooked := &hooked{item}
    this.list = append(this.list, hooked)

    return func() { this.remove(hooked) }
}

func (this *Hook) remove(hooked *hooked) {
    for i := range this.list {
        if this.list[i] == hooked {
            this.list = append(this.list[:i], this.list[i+1:]...)
            break
        }
    }
}

func (this *Hook) Call(arg1 io.Writer, arg2 ...interface{}) {
    for _, hooked := range this.list {
        hooked.Call(arg1, arg2...)
    }
}
`,
			wantErr: false,
		},
		{
			name: "Template{ Safe: true } uses sync.Mutex",
			this: HookTemplate{
				Imports: []string{
					"io",
				},
				InterfaceName: "Callback",
				PackageName:   "cbhook",
				MethodName:    "Call",
				MethodDeclArgs: []string{
					"arg1 io.Writer",
					"arg2 ...interface{}",
				},
				MethodCallArgs: []string{
					"arg1",
					"arg2...",
				},
				Safe: true,
			},
			// Use only spaces for indentation
			wantW: `// Code generated by hookgen; DO NOT EDIT.
// github.com/axard/things/cmd/hookgen

package cbhook

import (
    "io"
    "sync"
)

type Hook struct {
    list []*hookedItem
    m sync.Mutex
}

type hooked struct {
    Callback
}

type Cancel = func()

func (this *Hook) Append(item Callback) Cancel {
    this.m.Lock()
    defer this.m.Unlock()

    hooked := &hooked{item}
    this.list = append(this.list, hooked)

    return func() { this.remove(hooked) }
}

func (this *Hook) remove(hooked *hooked) {
    this.m.Lock()
    defer this.m.Unlock()

    for i := range this.list {
        if this.list[i] == hooked {
            this.list = append(this.list[:i], this.list[i+1:]...)
            break
        }
    }
}

func (this *Hook) Call(arg1 io.Writer, arg2 ...interface{}) {
    this.m.Lock()
    defer this.m.Unlock()

    for _, hooked := range this.list {
        hooked.Call(arg1, arg2...)
    }
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := tt.this.Write(w); (err != nil) != tt.wantErr {
				t.Errorf("Template.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Template.Write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
