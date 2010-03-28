package SortedList

import (
	//"sync"
	"container/vector"
	s "ghthor/Server"
)

type IntNodePtrServer struct {
	Server
	get chan (chan *IntNode)
	set chan *IntNode
	ptr *IntNode
}

func (ps *IntNodePtrServer) Init() {
	ps.Server.Init()
	ps.get = make(chan (chan *IntNode))
	// By making this channel a buffer we can clear
	// the buffer everytime there is a req to get
	ps.set = make(chan *IntNode, 10)
}

func (ps *IntNodePtrServer) Start() {
	if !ps.Running {
		ps.Running = true
		go func() {
			get := make(chan *IntNode)
			for {
				select {
					case <-ps.StopChan:
						ps.Running = false
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

func (is *IntNodeIndexSvr) Init() {
	is.Server.Init()
	is.get = make(chan (chan int64))
	is.offsetBy = make(chan int64, 200)
}

func (is *IntNodeIndexSvr) Start() {
	if !is.Running {
		is.Running = true
		go func() {
			get := make(chan int64)
			for {
				select {
					case <-is.StopChan:
						is.Running = false
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
	newIntNode.Init()
	return newIntNode
}

func (i *IntNode) Init() {
	i.WorkerPool.Init()

	i.next = new(IntNodePtrServer)
	i.next.Init()
	i.pushWorker(i.next)

	i.prev = new(IntNodePtrServer)
	i.prev.Init()
	i.pushWorker(i.prev)

	i.index = new(IntNodeIndexSvr)
	i.index.Init()
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
	getAndOffset chan (chan uint32)
	offsetBy chan int
}

func (ss *ListSizeSvr) Init() {
	ss.Server.Init()
	ss.get = make(chan (chan uint32))
	ss.getAndOffset = make(chan (chan uint32))
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
	if !ss.Running {
		ss.Running = true
		go func() {
			get := make(chan uint32)
			select {
				case <-ss.StopChan:
					ss.Running = false
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
				// This Implies a atomic read and then write
				case ss.getAndOffset <- get:
					//Empty the ss.offsetBy buffer
					offset, ok := <-ss.offsetBy
					for ok {
						ss.changeByOffset(offset)
						offset, ok = <-ss.offsetBy
					}
					get <- ss.size
					offset = <-ss.offsetBy
					ss.changeByOffset(offset)
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

	//goInsert (*SortedIntList) func(int)

	Insert chan int
	//Contains chan interface {}
	size *ListSizeSvr
}

func NewSortedIntList() *SortedIntList {
	newList := new(SortedIntList)
	newList.Init()
	return newList
}

func (sl *SortedIntList) Init() {
	sl.Server.Init()
	sl.WorkerPool.Init()

	sl.size = new(ListSizeSvr)
	sl.size.Init()
	sl.pushWorker(sl.size)
	sl.size.Start()

	sl.Insert = make(chan int)
	//sl.goInsert = &sl.goInsertInitial

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
	if !sl.Running {
		sl.Running = true
		go func() {
			for {
				select {
					case <-sl.StopChan:
						sl.Running = false
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
		sl.StopChan <- 0
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

func (sl *SortedIntList) goInsertNorm(val int) {
}

//func (sl *SortedIntList) goInsertInitial(val int) {
func (sl *SortedIntList) goInsert(val int) {
	// We Need and Atomic Read and Write here
	// Since multiple thread could attempt to insert the "first" element to the list
	// since we atomicly lock this for a getAndOffset all thread will lock until the size equal 1
	sizeLock := <-sl.size.getAndOffset
	curSize := <-sizeLock
	if curSize > 0 {
		// We have to release the lock, so that other threads can enter this section
		// Honestly I should do this with a state and change the value of this function
		sl.size.offsetBy <- 0
		//all calls to goInsert after this will be to goInsertNorm
		//sl.goInsert = sl.goInsertNorm
		sl.goInsertNorm(val)
	} else {
		sl.size.offsetBy <- 1
		insertee := newIntNode(val, sl)
		sl.insertHelper(sl.begin, insertee, sl.end)
		sl.middle = insertee
	}
}
