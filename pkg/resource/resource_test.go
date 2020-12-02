package resource

import (
	"testing"
)

func TestPackage(t *testing.T) {
	type args struct {
		resource string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Package() returns valid path if resource is full and valid",
			args: args{
				resource: "/path/to/package.Resource",
			},
			want: "/path/to/package",
		},
		{
			name: "Package() returns valid path if resource has only delimeter dot without object name",
			args: args{
				resource: "/path/to/package.",
			},
			want: "/path/to/package",
		},
		{
			name: "Package() returns valid path if resource hasn't object",
			args: args{
				resource: "/path/to/package",
			},
			want: "/path/to/package",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Package(tt.args.resource); got != tt.want {
				t.Errorf("Package() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObject(t *testing.T) {
	type args struct {
		resource string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Object() returns valid object name if resource is full and valid",
			args: args{
				resource: "/path/to/package.Resource",
			},
			want: "Resource",
		},
		{
			name: "Object() returns empty string if resource has only delimeter dot without object name",
			args: args{
				resource: "/path/to/package.",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				resource: "/path/to/package",
			},
			want: "",
		},
		{
			name: "",
			args: args{
				resource: ".Resource",
			},
			want: "Resource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Object(tt.args.resource); got != tt.want {
				t.Errorf("Object() = %v, want %v", got, tt.want)
			}
		})
	}
}
