package object

import (
	"container/vector"
)

type ObId uint64

type Object interface {
	Dispose(chan string)
}

type BaseObject struct {
	Id ObId
}

func (o *BaseObject) Init(interface {}) (*BaseObject) {
	return o
}

func (o *BaseObject) Dispose(disFinished chan string) {
}

//TODO: Implement the global object monitor on Node's!

//TODO: Record Object status's to a file to be used for optimization
// Like avg num of objects created so we can scale the objects vector
type ObjectList struct {
	objects vector.Vector
	numObs uint64
}

var ObjectKernel *ObjectList

func init() {
	ObjectKernel = new(ObjectList)
}
