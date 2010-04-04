package comm

import (
	//"fmt"
	i "ghthor/init"
	//"runtime"
)

// (type)Init structs Allow for Simplified Construction
// especially with complex embedded types
type CommInit struct {
	BufferLen int
}

type ComplexCommInit struct {
	CommInit
	NumChan int
}

// Defines a channel with embedded Buffer Length for easy cloning of an instance
type Comm struct {
	i.InitVar
	Comm chan interface {}
}

// Constructor
func (c *Comm) Init(arg interface {}) (interface {}) {
	// Call to embedded Init func that Checks for nil and reassigns arg to i.Default
	arg = c.InitVar.Init(arg)
	var initArg CommInit

/*
	// Allow for Weird combos of arg for initialization
	if x, ok := arg.(CommInit); ok {
		initArg = x
	} else if x, ok := arg.(int); ok {
		initArg = CommInit{BufferLen:x}
	} else {
		initArg = CommInit{BufferLen:0}
	}
*/

	switch arg.(type) {
		case *CommInit:
			initArg = *arg.(*CommInit)
		case CommInit:
			initArg = arg.(CommInit)
		case int:
			initArg = CommInit{BufferLen:arg.(int)}
		default:
			initArg = CommInit{BufferLen:0}
	}


	// Initialize c *Comm
	c.Comm = make(chan interface {}, initArg.BufferLen)
	c.InitArg = initArg
	return c
}

// Clean-up (Destructor)
func (c *Comm) Dispose(disposeComplete chan string) {
	//TODO: Maybe have this do some reporting, like if there are things left on the channel when it got closed
	close(c.Comm)
}

// These 2 struct's and func's Make the code a little ezier to read down the road
type CommIn struct {
	In chan interface {}
}
type CommOut struct {
	Out chan interface {}
}

func (c *Comm) AsIn() (in *CommIn) {
	in = new(CommIn)
	in.In = c.Comm
	return
}

func (c *Comm) AsOut() (out *CommOut) {
	out = new(CommOut)
	out.Out = c.Comm
	return
}

// Comm that has multiple channels of communication
type ComplexComm struct {
	i.InitVar
	Comm []*Comm
}

// Constructor
func (cc *ComplexComm) Init(arg interface {}) (interface {}) {
	// nil Check
	arg = cc.InitVar.Init(arg)
	var initArg ComplexCommInit
/*
	if x, ok := arg.(ComplexCommInit); ok {
		initArg = x
	} else if x, ok := arg.(int); ok {
		initArg = ComplexCommInit{NumChan:x}
	} else if x, ok := arg.(CommInit); ok {
		initArg = ComplexCommInit{NumChan:2, CommInit:x}
	} else {
		initArg = ComplexCommInit{NumChan:2} // CommInit.BufferLen defaults to 0
	}
*/

	switch arg.(type) {
		case *ComplexCommInit:
			initArg = *arg.(*ComplexCommInit)
		case ComplexCommInit:
			initArg = arg.(ComplexCommInit)
		case int:
			initArg.NumChan = arg.(int)
			initArg.BufferLen = 0
		case CommInit:
			initArg.NumChan = 2
			initArg.BufferLen = arg.(CommInit).BufferLen
		default:
			initArg.NumChan = 2
			initArg.BufferLen = 0
	}

	cc.Comm = make([]*Comm, initArg.NumChan)
	for i := 0; i < len(cc.Comm); i++ {
		cc.Comm[i] = new(Comm).Init(initArg.CommInit).(*Comm)
	}
	cc.InitArg = &initArg
	return cc
}

// Destructor
func (cc *ComplexComm) Dispose(disposeComplete chan string) {
	for i := 0; i < len(cc.Comm); i++ {
		cc.Comm[i].Dispose(disposeComplete)
	}
}

