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

type Node struct {
	val interface {}
	setVal chan int
	getVal chan int
	parent *SortedList
	next *Node
	prev *Node
	begin *Node
	end *Node
	index uint32
	lock sync.Mutex

	status string
	getStatus chan string

	stop chan bool
	setNext chan *Node
	setPrev chan *Node
}

func NewNode(val interface {}, parent *SortedList) (newNode *Node) {
	newNode = new(Node)
	newNode.val = val
	newNode.init(parent)
}

func (i *Node) init(parent *SortedList) {
	i.parent = parent
	i.stop = make(chan bool)
	i.getVal = make(chan int)
	i.setVal = make(chan int)
	i.getStatus = make(chan int)
	i.setNext = make(chan *Node)
	i.setPrev = make(chan *Node)
}

func (i *Node) Start() {
	go func() {
		i.status = "started"
		for {
			i.status = "paused"
			select {
				case i.getStatus <- i.status:
				case <- i.stop:
					i.status = "stopped"
					return
				case next := <-i.setNext:
					i.status = "working"
					//next.prev = i
					i.next = next
				case prev :=  <-i.setPrev:
					i.status = "working"
					//prev.next = i
					i.prev = prev
			}
		}
	}()
}

func (i *Node) SetVal(val int) {
	go func() {
		i.setVal <- val
	}()
}

func (i *Node) GetVal() (ret chan int) {
	ret = make(chan int)
	go func() {
		temp := <-getVal
		ret <- temp
	}()
	return ret
}

func (i *Node) SetNext(next *Node) {
	go func() {
		i.setNext <- next
	}()
}

func (i *Node) SetPrev(prev *Node) {
	go func() {
		i.setPrev <- prev
	}()
}

func (i *Node) Stop() {
	go func() {
		i.stop <- 0
	}()
}

type SortedList struct {
	WorkerPool
	begin *Node
	end *Node
	len int
}

func (il *SortedList) Insert(val int) {
	if il.begin == nil {
		il.begin = NewNode(val, il)
		il.PushWorker(il.begin)
	} else if il.end == nil {
		il.end = NewNode(val, il)
		il.PushWorker(il.end)
		il.begin.SetNext(il.end)
		il.end.SetPrev(il.begin)
	}
}

