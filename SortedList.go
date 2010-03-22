package SortedList

import (
	"sync"
)

type Worker interface {
	Start()
	Stop()
}

type WorkerPool struct {
	size int
	workers []Worker
}

func (wp *WorkerPool) init() {
	workers = make([]Worker, size)
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

func (i *IntNode) init() {
	stop = make(chan int)
	setNext = make(chan *IntNode)
	setPrev = make(chan *IntNode)
}

func (i *IntNode) Start() {
	go func() {
		for {
			select {
				case <- i.stop:
					return
				case next := <-i.setNext:
					next.prev = i
					i.next = next
				case prev :=  <-i.setPrev:
					prev.next = i
					i.prev = prev
			}
		}
	}()
}


func (i *IntNode) init(parent *IntList) bool {
	i.parent:
}

type IntList interface {
	At(index uint32) (*IntNode)
}
