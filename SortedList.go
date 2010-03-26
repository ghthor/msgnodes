package SortedList

import (
	"sync"
	"container/vector"
)

type Worker interface {
	Start()
	Stop()
}

type WorkerPool struct {
	size int
	//workers []Worker
	workers vector.Vector
}

func (wp *WorkerPool) init() {
	//wp.workers = make([]Worker, wp.size)
	wp.size = 100
}

func (wp *WorkerPool) PushWorker(worker Worker) {
	wp.workers.Push(worker)
	worker.Start()
}

func (wp *WorkerPool) Start() {
}

func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.(Worker).Stop()
	}
}

type IntNode struct {
	val int
	setVal chan int
	getVal chan int
	parent *SortedIntList
	next *IntNode
	prev *IntNode
	begin *IntNode
	end *IntNode
	index uint32
	lock sync.Mutex

	status string
	getStatus chan string

	stop chan int
	valServerStop chan int
	getNext chan *IntNode
	getPrev chan *IntNode
	setNext chan *IntNode
	setPrev chan *IntNode
}

func newIntNode(val int, parent *SortedIntList) (newIntNode *IntNode) {
	newIntNode = new(IntNode)
	newIntNode.val = val
	newIntNode.init(parent)
	return newIntNode
}

func (i *IntNode) init(parent *SortedIntList) {
	i.parent = parent
	i.parent.PushWorker(i)
	i.stop = make(chan int)
	i.valServerStop = make(chan int)
	i.getVal = make(chan int)
	i.setVal = make(chan int)
	i.getStatus = make(chan string)
	i.getNext = make(chan *IntNode)
	i.getPrev = make(chan *IntNode)
	i.setNext = make(chan *IntNode)
	i.setPrev = make(chan *IntNode)
}

func (i *IntNode) Start() {
	// This Loop emulates Locking via a select
	// On the channels that access this IntNode
	// In Theory This should work
	// The cases I'm wondering about are like when both a req for the prev pointer and a req to set th prev
	// Happen at the same time.  I think it will work but I haven't tested it yet
	go func() {
		i.status = "started"
		for {
			//i.status = "paused"
			select {
				case i.getStatus <- i.status:
				case <- i.stop:
					i.valServerStop <- 0
					i.status = "stopped"
					return
				// Set the next Pointer
				case next := <-i.setNext:
					//i.status = "working"
					i.next = next
				// Set the prev Pointer
				case prev := <-i.setPrev:
					//i.status = "working"
					i.prev = prev
				// Return the next Pointer
				case i.getNext <- i.next:
				// Return the prev Pointer 
				case i.getPrev <- i.prev:
			}
		}
	}()
	go func() {
		for {
			select {
				// Stop this go routine
				case <- i.valServerStop:
				// Return Val
				case i.getVal <- i.val:
				// Set Val
				case val := <-i.setVal:
					i.val = val
			}
		}
	}()
}

func (i* IntNode) Stop() {
	go func() {
		i.stop <- 0
	}()
}

type SortedIntList struct {
	WorkerPool
	begin *IntNode
	end *IntNode

	Insert chan int
	//Contains chan interface {}
	stop chan int
	sizeStop chan int
	size int
	//getSize chan (chan int)
	setSize chan int
	getSize chan int
}

func NewSortedIntList() *SortedIntList {
	newList := new(SortedIntList)

	// ---- Initializations ---- //
	newList.Insert = make(chan int)

	// The channel that kills the server processes
	newList.stop = make(chan int)

	// Channels that wrap up size
	newList.sizeStop = make(chan int)
	newList.size = 0
	//newList.getSize = make(chan (chan int))
	newList.getSize = make(chan int)
	newList.setSize = make(chan int)

	// Setup the Begin and End Nodes
	newList.begin = newIntNode(0, newList)
	newList.end = newIntNode(0, newList)
	newList.begin.setNext <- newList.end
	newList.end.setPrev <- newList.begin
	newList.PushWorker(newList)
	return newList
}

func (sl *SortedIntList) Start() {
	// Server for passing a value to insert or search for
	go func() {
		for {
			select {
				case <-sl.stop:
					sl.sizeStop <- 0
					return
				case newVal :=  <-sl.Insert:
					go sl.goInsert(newVal)
			}
		}
	}()
	// Server for reading and writing to the size of the List
	go func() {
		//setSize := make(chan int)
		for {
			select {
				case <-sl.sizeStop:
					return
				case newSize := <-sl.setSize:
					sl.size += newSize
				case sl.getSize <- sl.size:
					//<-setSize
					//sl.size = <-setSize
			}
		}
	}()
}

func (sl *SortedIntList) Stop() {
	go func() {
		sl.stop <- 0
	}()
}

func (sl *SortedIntList) insertHelper(prev *IntNode, insertee *IntNode,  next *IntNode) {
	insertee.next = next
	insertee.prev = prev
	go func() {
		prev.setNext <- insertee
	}()
	go func() {
		next.setPrev <- insertee
	}()
}

func (sl *SortedIntList) goInsert(val int) {
	curSize := <-sl.getSize
	if curSize != 0 {
		//sizeLock <- sl.size
	} else {
		//sl.size++
		sl.setSize <- +1
		//sizeLock <- sl.size
		insertee := newIntNode(val, sl)
		sl.insertHelper(sl.begin, insertee, sl.end)
	}
}
