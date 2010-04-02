package node

import (
	"fmt"
	i "ghthor/init"
	c "ghthor/comm"
	//"runtime"
)

// Struct embeds dir prop in the channel
// This Struct is for abstracting the NodeConn Struct into something more readable
type NodeComm struct {
	CommIn
	CommOut
}

type BaseNode struct {
	i.InitVar
	MsgProcessor
	NodeComm
	shutDown chan int
}

type BaseNodeInit struct {
	in *c.Comm
	out *c.Comm
	shutDown chan int
}

func (n *BaseNode) Init(arg interface {}) (interface {}) {
	arg = n.InitVar.Init(arg)
	var initArg *BaseNodeInit
	switch arg.(type) {
		case BaseNodeInit:
			initArg = &arg.(BaseNodeInit)
		case *BaseNodeInit:
			initArg = arg.(*BaseNodeInit)
		case c.Comm:
			temp := arg.(c.Comm)
			initArg = &BaseNodeInit{in:&c.Comm, out:&c.Comm}
		case *c.Comm:
			temp := arg.(*c.Comm)
			initArg = &BaseNodeInit{in:c.Comm, out:c.Comm}
		default:
			//TODO: Invalid Initializer
			//Really difficult to have a base case becuase I don't want the nodes to be creating Comm objects
			// The reason they can't create comm objects is because they Don't ever have a pointer to a comm object
			return
	}
	n.InitArg = initArg
	n.in = initArg.in.AsIn()
	n.out = initArg.out.AsOut()
	if initArg.shutDown != nil {
		n.shutDown = initArg.shutDown
	} else {
		n.shutDown = make(chan int, 1)
	}

	return n
}

func (n *BaseNode) Listen() {
}

//func (n *BaseNode) OpenProxy(

type Msg struct {
	priority uint32
	propagate bool
	status chan string
	str string
}

type ShutdownMsg struct {
	Msg
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

