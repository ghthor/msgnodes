package comm

import (
	//"fmt"
	i "ghthor/init"
	//"runtime"
)

type CommInit struct {
	BufferLen int
}

type ComplexCommInit struct {
	CommInit
	NumChan int
}

type Comm struct {
	i.InitVar
	comm chan interface {}
}

func (c *Comm) Init(arg interface {}) (interface {}) {
	arg = c.InitVar.Init(arg)
	var initArg CommInit

	// Allow for Weird combos of arg for initialization
	switch arg.(type) {
		//case *CommInit:
			//initArg = arg.(*CommInit)
		case CommInit:
			initArg = arg.(CommInit)
		case int:
			initArg = CommInit{BufferLen:arg.(int)}
		// May Implement this Differently
		// - Could be used like cloned := clonee.Init(init.Clone{}).(*typeOfClonee) or cloned.Init(clonee)
		case i.Clone:
		case i.Default:
			initArg = CommInit{BufferLen:0}
		default:
			// TODO: make this throw some kinda error
			initArg = CommInit{BufferLen:0}
	}
	if initArg.BufferLen == 0 {
		c.comm = make(chan interface {})
	} else {
		c.comm = make(chan interface {}, initArg.BufferLen)
	}
	c.InitArg = &initArg
	return c
}

func (c *Comm) Dispose(disposeComplete chan string) {
	//TODO: Maybe have this do some reporting, like if there are things left on the channel when it got closed
	close(c.comm)
}

// These 2 struct's Makes the code a little ezier to read down the road
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

func (cc *ComplexComm) Init(arg interface {}) (interface {}) {
	arg = cc.InitVar.Init(arg)
	var initArg ComplexCommInit
	switch arg.(type) {
		case ComplexCommInit:
			initArg = arg.(ComplexCommInit)
		case int:
			initArg.NumChan = arg.(int)
			initArg.BufferLen = 0
		case CommInit:
			initArg.NumChan = 2
			initArg.BufferLen = arg.(CommInit).BufferLen
		case i.Default:
			initArg.NumChan = 2
			initArg.BufferLen = 0
		default:
			//TODO: Make this throw an error
			initArg.NumChan = 2
			initArg.BufferLen = 0
	}
	cc.comm = make([]*Comm, initArg.NumChan)
	for i := 0; i < len(cc.comm); i++ {
		cc.comm[i] = new(Comm).Init(initArg.CommInit).(*Comm)
	}
	cc.InitArg = &initArg
	return cc
}

func (cc *ComplexComm) Dispose(disposeComplete chan string) {
	for i := 0; i < len(cc.comm); i++ {
		cc.comm[i].Dispose(disposeComplete)
	}
}

type ListComm struct {
}

//func (lc *ListComm)

// Struct embeds dir prop in the channel
// This Struct is for abstracting the NodeConn Struct into something more readable
type NodeComm struct {
	CommIn
	CommOut
	dir string
}

// This is the class that all nodes communicate through
// TODO: buffer and sort msg's by priority
type NodeConn struct {
	next chan interface {}
	prev chan interface {}
}

// Create the channels for communicating
func (n *NodeConn) Init() (interface {}) {
	n.next = make(chan interface {})
	n.prev = make(chan interface {})
	return n
}

func (n *NodeConn) Close() {
	close(n.next)
	close(n.prev)
}

// abstract this conn into and in/out NodeComm to the next Node
func (n *NodeConn) GetAsNext() (next *NodeComm) {
	next = new(NodeComm)
	next.in = n.prev
	next.out = n.next
	next.dir = "Next->"
	return next
}

// abstract this conn into and in/out NodeComm to the prev Node
func (n *NodeConn) GetAsPrev() (prev *NodeComm) {
	prev = new(NodeComm)
	prev.in = n.next
	prev.out = n.prev
	prev.dir = "<-Prev"
	return prev
}

