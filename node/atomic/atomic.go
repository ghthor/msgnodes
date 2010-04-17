package atomic_int64

import (
	"ghthor/node"
	"ghthor/node/buffer"
)

type Msg node.Msg
type BaseMsg node.BaseMsg

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
	buffer.BufferNode
	val Type
}

func (an *AtomicNode) Init(val Type, bufferSz int, ShutDown chan int) (*AtomicNode) {
	an.val = val
	an.BufferNode.Init(bufferSz, ShutDown)
	return an
}
