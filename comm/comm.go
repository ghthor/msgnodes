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
	comm chan interface {}
}

// Constructor
func (c *Comm) Init(arg interface {}) (interface {}) {
	// Call to embedded Init func that Checks for nil and reassigns arg to i.Default
	arg = c.InitVar.Init(arg)
	var initArg CommInit

	// Allow for Weird combos of arg for initialization
	if x, ok := arg.(CommInit); ok {
		initArg = x
	} else if x, ok := arg.(int); ok {
		initArg = CommInit{BufferLen:x}
	} else {
		initArg = CommInit{BufferLen:0}
	}

	// Initialize c *Comm
	c.comm = make(chan interface {}, initArg.BufferLen)
	c.InitArg = &initArg
	return c
}

// Clean-up (Destructor)
func (c *Comm) Dispose(disposeComplete chan string) {
	//TODO: Maybe have this do some reporting, like if there are things left on the channel when it got closed
	close(c.comm)
}

// These 2 struct's and func's Make the code a little ezier to read down the road
type CommIn struct {
	in chan interface {}
}
type CommOut struct {
	out chan interface {}
}

func (c *Comm) AsIn() (in *CommIn) {
	in = new(CommIn)
	in.in = c.comm
	return
}

func (c *Comm) AsOut() (out *CommOut) {
	out = new(CommOut)
	out.out = c.comm
	return
}

// Comm that has multiple channels of communication
type ComplexComm struct {
	i.InitVar
	comm []*Comm
}

// Constructor
func (cc *ComplexComm) Init(arg interface {}) (interface {}) {
	// nil Check
	arg = cc.InitVar.Init(arg)
	var initArg ComplexCommInit
	if x, ok := arg.(ComplexCommInit); ok {
		initArg = x
	} else if x, ok := arg.(int); ok {
		initArg = ComplexCommInit{NumChan:x}
	} else if x, ok := arg.(CommInit); ok {
		initArg = ComplexCommInit{NumChan:2, CommInit:x}
	} else {
		initArg = ComplexCommInit{NumChan:2} // CommInit.BufferLen defaults to 0
	}

/*	What I thought was that the statement
		case ComplextCommInit:
	was taking the place of
		if x, ok := arg.(ComplexCommInit); ok
	But now I understand that
		switch arg.(type)
	is just querying arg's declared type which is interface {}
	I misunderstodd Type Switch, for Type Assertion Switch

	switch arg.(type) {
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
*/
	cc.comm = make([]*Comm, initArg.NumChan)
	for i := 0; i < len(cc.comm); i++ {
		cc.comm[i] = new(Comm).Init(initArg.CommInit).(*Comm)
	}
	cc.InitArg = &initArg
	return cc
}

// Destructor
func (cc *ComplexComm) Dispose(disposeComplete chan string) {
	for i := 0; i < len(cc.comm); i++ {
		cc.comm[i].Dispose(disposeComplete)
	}
}

