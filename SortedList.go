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
	workers container.Vector
}

func (wp *WorkerPool) init() {
	//wp.workers = make([]Worker, wp.size)
	size = 100
	wp.workers = new(container.Vector)
}

func (wp *WorkerPool) PushWorker(worker *Worker) {
	wp.workers
}

func (wp *WorkerPool) Start() {
}

func (wp *WorkerPool) Stop() {
}

type IntNode struct {
	val int
	parent *IntList
	next *IntNode
	prev *IntNode
	begin *IntNode
	end *IntNode
	index uint32
	lock sync.Mutex

	stop chan int
	setNext chan *IntNode
	setPrev chan *IntNode
}

func (i *IntNode) init(parent *IntList) {
	i.parent = parent
	i.stop = make(chan int)
	i.setNext = make(chan *IntNode)
	i.setPrev = make(chan *IntNode)
}

func (i *IntNode) Start() {
	go func() {
		for {
			select {
				case <- i.stop:
					return
				case next := <-i.setNext:
					//next.prev = i
					i.next = next
				case prev :=  <-i.setPrev:
					//prev.next = i
					i.prev = prev
			}
		}
	}()
}

type IntList interface {
	At(index uint32) (*IntNode)
}
