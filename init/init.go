package init

type Initialize interface {
	Init(interface {}) (interface {})
}

type InitVar struct {
	InitArg interface {}
}

type Default struct {
}

// Used to check for a nil value passed to a Init() call
// If nil is passed this turns arg into a Default
func (iv *InitVar) Init(arg interface {}) (interface {}) {
	if arg == nil {
		arg = Default{}
	}
	return arg
}

type Clone struct {
}
