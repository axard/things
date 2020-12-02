package resource

import (
	"path"
	"strings"
)

func Package(resource string) string {
	d, f := path.Split(resource)
	return d + strings.TrimSuffix(f, path.Ext(f))
}

func Object(resource string) string {
	return strings.TrimLeft(path.Ext(resource), ".")
}
