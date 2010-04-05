package atomic_int64

import (
	"ghthor/node"
)

type AtomicNode struct {
	node.BaseNode
	val Type
}

func (an *AtomicNode) Init(val Type, ShutDown chan int) (*AtomicNode) {
	an.val = val
	an.BaseNode.Init(ShutDown)
	return an
}
