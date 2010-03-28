package SortedList

import (
	//"sync"
	"container/vector"
)

type Worker interface {
	Start()
	//Pause()
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

func (wp *WorkerPool) pushWorker(worker Worker) {
	wp.workers.Push(worker)
	//worker.Start()
}

func (wp *WorkerPool) Start() {
	for _, worker := range wp.workers {
		worker.(Worker).Start()
	}
}

func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.(Worker).Stop()
	}
}

type Server struct {
	running bool
	//start chan int
	//pause chan int
	stop chan int
}

func (s *Server) init() {
	s.running = false
	//s.start = make(chan int)
	//s.pause = make(chan int)
	s.stop = make(chan int)
}

func (s *Server) Start() {
	//if !s.running {
		//go func() {
			//s.start <- 0
		//}()
	//}
}

func (s *Server) Stop() {
	if s.running {
		go func() {
			s.stop <- 0
		}()
	}
}

type IntNodePtrServer struct {
	Server
	get chan (chan *IntNode)
	set chan *IntNode
	ptr *IntNode
}

func (ps *IntNodePtrServer) init() {
	ps.Server.init()
	ps.get = make(chan (chan *IntNode))
	// By making this channel a buffer we can clear
	// the buffer everytime there is a req to get
	ps.set = make(chan *IntNode, 10)
}

func (ps *IntNodePtrServer) Start() {
	if !ps.running {
		ps.running = true
		go func() {
			get := make(chan *IntNode)
			for {
				select {
					case <-ps.stop:
						ps.running = false
						return
					case ptr := <-ps.set:
						ps.ptr = ptr
					case ps.get <- get:
						// Clear the Buffer to set
						ptr, ok := <-ps.set
						for ok {
							ps.ptr = ptr
							ptr, ok = <-ps.set
						}
						get <- ps.ptr
				}
			}
		}()
	}
}

type IntNodeIndexSvr struct {
	Server
	val int64
	get chan (chan int64)
	offsetBy chan int64
}

func (is *IntNodeIndexSvr) init() {
	is.Server.init()
	is.get = make(chan (chan int64))
	is.offsetBy = make(chan int64, 200)
}

func (is *IntNodeIndexSvr) Start() {
	if !is.running {
		is.running = true
		go func() {
			get := make(chan int64)
			for {
				select {
					case <-is.stop:
						is.running = false
						return
					case offset := <-is.offsetBy:
						is.val += offset
					case is.get <- get:
						//empty the offsetBy Buffer
						offset, ok := <-is.offsetBy
						for ok {
							is.val += offset
							offset, ok = <-is.offsetBy
						}
						get <- is.val
				}
			}
		}()
	}
}

type IntNode struct {
	WorkerPool
	//Server
	val int
	//setVal chan int
	//getVal chan int
	parent *SortedIntList
	next *IntNodePtrServer
	prev *IntNodePtrServer
	index *IntNodeIndexSvr
}

func newIntNode(val int, parent *SortedIntList) (newIntNode *IntNode) {
	newIntNode = new(IntNode)
	newIntNode.val = val
	newIntNode.parent = parent
	//newIntNode.parent.PushWorker(newIntNode)
	newIntNode.init()
	return newIntNode
}

func (i *IntNode) init() {
	i.WorkerPool.init()

	i.next = new(IntNodePtrServer)
	i.next.init()
	i.pushWorker(i.next)

	i.prev = new(IntNodePtrServer)
	i.prev.init()
	i.pushWorker(i.prev)

	i.index = new(IntNodeIndexSvr)
	i.index.init()
	i.pushWorker(i.index)
	// Lockers for accessing the value of the element
	//i.getVal = make(chan int)
	//i.setVal = make(chan int)
}

func (i *IntNode) Start() {
	i.WorkerPool.Start()
}

func (i* IntNode) Stop() {
	i.WorkerPool.Stop()
}

type ListSizeSvr struct {
	Server
	size uint32
	get chan (chan uint32)
	offsetBy chan int
}

func (ss *ListSizeSvr) init() {
	ss.Server.init()
	ss.get = make(chan (chan uint32))
	ss.offsetBy = make(chan int, 200)
}

func (ss *ListSizeSvr) changeByOffset(offset int) {
	if offset == 1 {
		ss.size += 1
	} else if offset == -1 {
		ss.size -= 1
	}
}

func (ss *ListSizeSvr) Start() {
	if !ss.running {
		ss.running = true
		go func() {
			get := make(chan uint32)
			select {
				case <-ss.stop:
					ss.running = false
					return
				case offset := <-ss.offsetBy:
					ss.changeByOffset(offset)
				case ss.get <- get:
					//Empty the ss.offsetBy buffer
					offset, ok := <-ss.offsetBy
					for ok {
						ss.changeByOffset(offset)
						offset, ok = <-ss.offsetBy
					}
					get <- ss.size
			}
		}()
	}
}

type SortedIntList struct {
	WorkerPool
	Server
	begin *IntNode
	end *IntNode
	middle *IntNode

	Insert chan int
	//Contains chan interface {}
	size *ListSizeSvr
}

func NewSortedIntList() *SortedIntList {
	newList := new(SortedIntList)
	newList.init()
	return newList
}

func (sl *SortedIntList) init() {
	sl.Server.init()
	sl.WorkerPool.init()

	sl.size = new(ListSizeSvr)
	sl.size.init()
	sl.pushWorker(sl.size)
	sl.size.Start()

	sl.Insert = make(chan int)

	sl.begin = newIntNode(0, sl)
	sl.pushWorker(sl.begin)
	sl.begin.Start()

	sl.end = newIntNode(0, sl)
	sl.pushWorker(sl.end)
	sl.end.Start()
	sl.end.index.val = 1
}

func (sl *SortedIntList) Start() {
	// Server for passing a value to insert or search for
	go sl.WorkerPool.Start()
	if !sl.running {
		sl.running = true
		go func() {
			for {
				select {
					case <-sl.stop:
						sl.running = false
						return
					case newVal :=  <-sl.Insert:
						go sl.goInsert(newVal)
				}
			}
		}()
	}
}

func (sl *SortedIntList) Stop() {
	go sl.WorkerPool.Stop()
	go func() {
		sl.stop <- 0
	}()
}

// this func is used once the location is discovered where the node goes
func (sl *SortedIntList) insertHelper(prev *IntNode, insertee *IntNode,  next *IntNode) {
	insertee.next.set <- next
	insertee.prev.set <- prev
	indexLock := <-prev.index.get
	prevIndex := <- indexLock
	insertee.index.offsetBy <- (prevIndex + 1)
	go func() {
		prev.next.set <- insertee
	}()
	go func() {
		next.prev.set <- insertee
	}()
}

func (sl *SortedIntList) goInsert(val int) {
	sizeLock := <-sl.size.get
	curSize := <-sizeLock
	if curSize != 0 {
	} else {
		sl.size.offsetBy <- 1
		insertee := newIntNode(val, sl)
		sl.insertHelper(sl.begin, insertee, sl.end)
		sl.middle = insertee
	}
}
