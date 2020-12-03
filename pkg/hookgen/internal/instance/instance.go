package instance

type (
	Struct struct{}

	Interface0 interface {
		Method(interface{})
	}

	Interface1 interface {
		Method(...interface{})
	}

	Interface2 interface {
		Method(s string, _ interface{})
	}

	Interface3 interface {
		Method(i int, s Struct)
	}
)
