package node

import (
	"fmt"
	//"runtime"
)


type Comm struct {
	comm chan interface {}
}

func (c *Comm) Init() (interface {}) {
	c.comm = make(chan interface {})
	return c
}

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

type DblComm struct {
	comm []*Comm
}

func (dc *DblComm) Init() (interface {}) {
	chan1 = (new(Comm)).Init()
	chan2

type ListComm struct {
	DblComm
}

func (lc *ListComm)

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

type SimpleNode struct {
	name string
	comm *NodeComm
	shutDown chan int
}

// A Node
type Node struct {
	name string
	prev *NodeComm
	next *NodeComm
	shutDown chan int
}

type ControlNode struct {
	Node
}

type NodeList struct {
	Node
	begin *NodeComm
	end	*NodeComm
	size uint64
}

func (n *Node) Init() (interface {}) {
	n.shutDown = make(chan int, 10)
	return n
}

type Msg struct {
	priority uint32
	propagate bool
	status chan string
	str string
}

type ShutdownMsg struct {
	Msg
	val string
	from *Node
	complete chan string
}

// Process a Msg
func (n *Node) process(msg interface {}) (outMsg interface {}, msgStr string) {
	switch msg.(type) {
		case string:
			outMsg = msg
			msgStr = msg.(string)
		case Msg:
			outMsg = msg
			msgStr = msg.(Msg).str
		case ShutdownMsg:
			sdMsg := msg.(ShutdownMsg)
			msgStr = fmt.Sprint("ShutdownMsg from: ", sdMsg.from.name)
			sdMsg.from = n
			outMsg = sdMsg
			n.shutDown <- 0
		default:
			outMsg = &Msg{propagate:false, str:"Unknown Msg"}
			msgStr = outMsg.(Msg).str
	}
	return outMsg, msgStr
}

// Testing Func, drops a msg into the chain
func (n *Node) DropMsg(msg interface {}, dir *NodeComm) {
	msg, msgStr := n.process(msg)
	ChanPrintln <- fmt.Sprint("Msg Dropped in ", n.name, ", msg: ", msgStr)
	go func() {
		dir.out <- msg
	}()
}

// Connect a Node to its neighboring NodeConn's represented as NodeComm's
func (n *Node) connect(prev *NodeConn, next *NodeConn) {
	if prev != nil {
		n.prev = prev.GetAsPrev()
	}
	if next != nil {
		n.next = next.GetAsNext()
	}
}

// End of a list of Node's Proxy
func (n *Node) openProxyEndPt(comm *NodeComm) {
	if comm != nil {
		go func() {
			defer func() { ChanPrintln <- fmt.Sprint(n.name, ": ", "EndPt comm closed\nIn From: ", comm.dir, "\n"); }()
			for !closed(comm.in) && !closed(comm.out) {
				select {
					case msg := <-comm.in:
						if msg == nil {
							break
						}
						msg, msgStr := n.process(msg)
						ChanPrintln <- fmt.Sprint(n.name, ": ", msgStr)
						sdMsg, isSD := msg.(ShutdownMsg)
						if isSD {
							sdMsg.complete <- n.name
						}
					case shutDown := <-n.shutDown:
						// Nonblocking since it is buffered, ensures that all other "server" go routines exit
						n.shutDown <- shutDown
						return
				}
			}
		}()
	}
}

// Proxy all msg's going through this node from in to out, doing processing inbetween
func (n *Node) openProxy(in *NodeComm, out *NodeComm) {
	if (in != nil) && (out != nil) {
		go func() {
			defer func() { ChanPrintln <- fmt.Sprint(n.name, ": ", "comm closed\nDir: ", out.dir, "\n"); }()
			for !closed(in.in) && !closed(out.out) {
				select {
					case msg := <-in.in:
						if msg == nil {
							break
						} else {
							outMsg, msgStr := n.process(msg)
							ChanPrintln <- fmt.Sprint(n.name, ": ", msgStr)
							out.out <- outMsg
						}
					case shutDown := <-n.shutDown:
						n.shutDown <- shutDown
						return
				}
			}
		}()
	}
}

func (n *Node) Start() {
}

