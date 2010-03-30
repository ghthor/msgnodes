package object

import (
	"container/vector"
	i "ghthor/init"
)

type ObId uint64

type Object interface {
	i.Initialize
	Dispose(chan string)
}

type BaseObject struct {
	i.InitVar
	Id ObId
}

func (o *BaseObject) Init(interface {}) (interface {}) {
	return o
}

func (o *BaseObject) Dispose(disFinished chan string) {
}

//TODO: finish this after constructing the comm Package
//type IdSvr struct {


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
