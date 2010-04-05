package node

import (
	//"fmt"
	//i "ghthor/init"
	//c "ghthor/comm"
	//"runtime"
	"sync"
)

type MsgPriority uint64

type Msg interface {
	Priority()(uint64) 
	SetRecvId(uint64)
	SetProcId(uint64)
	RecvId() (uint64)
	ProcId() (uint64)
}

type BaseMsg struct {
	recvId uint64
	procId uint64
}

func (m *BaseMsg) SetRecvId(id uint64) { m.recvId = id }
func (m *BaseMsg) RecvId() (uint64) { return m.recvId }
func (m *BaseMsg) SetProcId(id uint64) { m.procId = id }
func (m *BaseMsg) ProcId() (uint64) { return m.procId }

type Node interface {
	PassMsg() (chan Msg)
	Stop()
	Listen()
}

type BaseNode struct {
	Running bool
	ShutDownCh chan int
	sync.Mutex
}

// Sets the ShutDownCh channel if it is going to linked with other Node's Shutdown channels
func (n *BaseNode) Init(ShutDownCh chan int) (*BaseNode) {
	n.Running = false
	if ShutDownCh == nil {
		n.ShutDownCh = make(chan int, 1)
	} else {
		n.ShutDownCh = ShutDownCh
	}
	return n
}

func (n *BaseNode) Dispose() {
	close(n.ShutDownCh)
}

func (n *BaseNode) Stop() {
	go func() {
		n.Lock()
		if n.Running {
			n.Unlock()
			n.ShutDownCh <- 0
			return
		}
		n.Unlock()
	}()
}

func (n *BaseNode) ShutDown(sdVal int) {
	n.Lock()
	n.Running = false
	n.Unlock()
	sdVal++
	n.ShutDownCh <- sdVal
}

/*
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
			n.ShutDownCh <- 0
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
					case ShutDownCh := <-n.ShutDownCh:
						// Nonblocking since it is buffered, ensures that all other "server" go routines exit
						n.ShutDownCh <- ShutDownCh
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
					case ShutDownCh := <-n.ShutDownCh:
						n.ShutDownCh <- ShutDownCh
						return
				}
			}
		}()
	}
}
*/
