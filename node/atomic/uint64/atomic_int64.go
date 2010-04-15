package atomic_int64

import (
	"ghthor/node"
	"ghthor/node/buffer"
)

type Msg node.Msg

type QueryVal struct {
	node.BaseMsg
	query chan Type
}

type SetTo struct {
	node.BaseMsg
	newVal Type
}

type QueryAndSet struct {
	QueryVal
}

// Only Needed For uint's!
type OffsetDir byte
const (
	Pos OffsetDir = iota;
	Neg
)

type OffsetBy struct {
	node.BaseMsg
	Offset Type
	Dir OffsetDir
}

type Monitor struct {
	node.BaseMsg
	comm chan Msg
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
