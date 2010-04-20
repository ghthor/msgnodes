package atomic_uint64

import (
	"ghthor/node"
	"ghthor/node/buffer"
)

type Msg node.Msg
type BaseMsg node.BaseMsg
type BaseNode node.BaseNode

type QueryVal struct {
	BaseMsg
	Query chan Type
}

type QueryAndSet struct {
	QueryVal
}

type SetTo struct {
	BaseMsg
	NewVal Type
}

type OffsetBy struct {
	BaseMsg
	Offset int
}

type Monitor struct {
	BaseMsg
	Comm chan Msg
}

type RemoveMonitor struct {
	Monitor *Monitor
}

type AtomicNode struct {
	BaseNode
	buffer.BufferNode
	val Type
}

func (an *AtomicNode) Init(val Type, bufferSz int, ShutDown chan int) (*AtomicNode) {
	an.val = val
	an.BufferNode.Init(bufferSz, ShutDown)
	return an
}

func (an *AtomicNode) Listen() {
	an.BufferNode.Listen()
	an.Lock()
	defer an.Unlock()
	if !an.Running {
		an.Running = true
		go func() {
			for !closed(an.In) || !closed(an.ShutDownCh) {
				select {
					case sdVal := <-an.ShutDownCh:
						an.ShutDown(sdVal)
						return
					case MsgChan := <-an.MsgReq:
						msg := <-MsgChan
						an.processMsg(msg)
				}
			}
		}()
	}
}
