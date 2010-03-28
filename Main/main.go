package main

import (
	"fmt"
)

// Struct embeds dir prop in the channel

// This Struct is for abstracting the NodeConn Struct into something more readable
type NodeComm struct {
	in chan interface {}
	out chan interface {}
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

// A Node
type Node struct {
	name string
	prev *NodeComm
	next *NodeComm
	shutDown chan int
}

func (n *Node) Init() (interface {}) {
	n.shutDown = make(chan int, 2)
	return n
}

type Msg struct {
	priority uint32
	propagate bool
	str string
}

//type UnknownMsg &Msg{propagate:false}

type ShutdownMsg struct {
	Msg
	val string
	from *Node
	complete chan string
}

// Process a Msg
// TODO: Make this do some Processing
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
func (n *Node) dropMsg(msg interface {}, dir *NodeComm) {
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
	//for (closed(prev.comm) && closed(next.comm)) || (prev == nil && next == nil) {
		//select {
			//case
}

var ChanPrintln chan string

func init() {
	ChanPrintln = make(chan string, 500)
}

func Dump() {
	for !closed(ChanPrintln) {
		msg := <-ChanPrintln
		if msg == "" {
			msg = "I think the channel is closed"
			close(ChanPrintln)
		}
		fmt.Println(msg)
	}
}

func main() {


	conn1 := (new(NodeConn).Init()).(*NodeConn)
	conn2 := (new(NodeConn).Init()).(*NodeConn)

	node1 := &Node{name:"Begin"}
	node2 := &Node{name:"Mid"}
	node3 := &Node{name:"End"}

	node1.Init()
	node2.Init()
	node3.Init()

	node1.connect(nil, conn1)
	node2.connect(conn1, conn2)
	node3.connect(conn2, nil)

	node1.openProxyEndPt(node1.next)
	node2.openProxy(node2.prev, node2.next)
	node2.openProxy(node2.next, node2.prev)
	node3.openProxyEndPt(node3.prev)

	//go func() {
	go func() {
		node1.dropMsg("Dropped in Node1", node1.next)
	}()

	go func() {
		node3.dropMsg("Dropped in Node3", node3.prev)
	}()

	go func() {
		msg := &ShutdownMsg{from:node1, complete:make(chan string)}
		msg.propagate = true
		node1.dropMsg(*msg, node1.next)
		nodeName := <-msg.complete
		ChanPrintln <- fmt.Sprint("Shutdown Completed, ended on >", nodeName)
		close(ChanPrintln)
	}()

	fmt.Printf("\n\n");

	Dump()
	conn1.Close()
	conn2.Close()

	fmt.Printf("\n\n");
}
