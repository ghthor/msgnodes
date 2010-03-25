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

type Comparable interface {
	LessThan(comp Comparable) bool
}

type CompInt int

func (ci *int) LessThan(comp Comparable) bool {
	return (ci < comp)
}

type Node struct {
	val Comparable
	setVal chan Comparable
	getVal chan Comparable
	parent *SortedList
	next *Node
	prev *Node
	begin *Node
	end *Node
	index uint32
	lock sync.Mutex

	status string
	getStatus chan string

	stop chan int
	valServerStop chan int
	getNext chan *Node
	getPrev chan *Node
	setNext chan *Node
	setPrev chan *Node
}

func NewNode(val interface {}, parent *SortedList) (newNode *Node) {
	newNode = new(Node)
	newNode.val = val
	newNode.init(parent)
}

func (i *Node) init(parent (*SortedList) {
	i.parent = parent
	i.stop = make(chan int)
	i.valServerStop = make(chan int)
	i.getVal = make(chan Comparable)
	i.setVal = make(chan Comparable)
	i.getStatus = make(chan string)
	i.getNext = make(chan *Node)
	i.getPrev = make(chan *Node)
	i.setNext = make(chan *Node)
	i.setPrev = make(chan *Node)
}

func (i *Node) Start() {
	// This Loop emulates Locking via a select
	// On the channels that access this Node
	// In Theory This should work
	// The cases I'm wondering about are like when both a req for the prev pointer and a req to set th prev
	// Happen at the same time.  I think it will work but I haven't tested it yet
	go func() {
		i.status = "started"
		for {
			i.status = "paused"
			select {
				case i.getStatus <- i.status:
				case <- i.stop:
					i.valServerStop <- 0
					i.status = "stopped"
					return
				// Set the next Pointer
				case next := <-i.setNext:
					i.status = "working"
					i.next = next
				// Set the prev Pointer
				case prev :=  <-i.setPrev:
					i.status = "working"
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
				case <- valServerStop:
				// Return Val
				case getVal <- val:
				// Set Val
				case val <- setVal:
			}
		}
	}()
}



type SortedList struct {
	WorkerPool
	begin *Node
	end *Node

	Insert chan Comparable
	//Contains chan interface {}
	stop chan int
	len int
}

func NewSortedList() *SortedList {
	newList := new(SortedList)
	newList.Insert = make(chan Comparable)
	newList.stop = make(chan int)
	newList.len = 0
	newList.PushWorker(newList)
	return newList
}

func (sl *SortedList) Start() {
	go func() {
		for {
			select {
				case exitVal <- sl.stop:
					return
				case newVal <- sl.Insert:
			}
		}
	}()
}

func (sl *SortedList) Stop() {
	go func() {
		sl.stop <- 0
	}()
}

