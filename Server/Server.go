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
	started vector.Vector
	stopped vector.Vector
	workers vector.Vector
}

func (wp *WorkerPool) Init() {
	wp.size = 0
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

StopAllSvr := make(chan uint64, 1)
ExitLog := make(chan string, 2000)

type FailsafeStop struct {
	ExitWaitingChan func()
	sync.Mutex
}

func (fs *FailsafeStop) SetFunc(stopFunc func()) {
	fs.Lock()
	defer fs.Unlock()
	fs.ExitWaitingChan = stopFunc
}

func (fs *FailsafeStop) Exec() {
	fs.Lock()
	defer fs.Unlock()
	if fs.ExitWaitingChan != nil {
		go fs.ExitWaitingChan()
	}
}

// Wrap up the Running Variable
type Running struct {
	running bool
	sync.RWMutex
}

func (r *Running) Get() bool {
	r.RLock()
	defer r.RUnlock()
	return running
}

func (r *Running) SetTo(running bool) {
	r.Lock()
	defer r.Unlock()
	r.running = running
}

// Returns true if running wasn't already ==  true
func (r *Running) ToggleOn() bool {
	r.Lock()
	defer r.Unlock()
	if !r.running {
		r.running = true
	} else {
		return false
	}
	return true
}

// Returns true if running wasn't already == false
func (r *Running) ToggleOff() bool {
	r.Lock()
	defer r.Unlock()
	if r.running {
		r.running = false
	} else {
		return false
	}
	return true
}

type Server struct {
	Running
	Id uint64
	StopAllSvr chan uint64
	ExitLog chan string
	StopChan chan int

	Failsafe FailsafeStop
}

func (s *Server) Init() {
	s.SetTo(false)
	s.StopAllSvr = StopAllSvr
	s.ExitLog = ExitLog
	s.StopChan = make(chan int)
}

func (s *Server) Start() {
	go func() {
		numStopped := <-s.StopAllSvr
		numStopped++
		s.StopAllSvr <- numStopped
	}()
}

func (s *Server) Stop() {
	s.FailSafe.Exec()
	if s.Running.Get() {
		go func() {
			s.StopChan <- 0
		}()
	}
}

// class MainSvr : WorkerPool, Server
type MainSvr struct {
	WorkerPool
	Server
}

// class SvrID : Server
type SvrID struct {
	Server
	totalSize uint64
	NewID chan uint64
}

func (s *SvrID) Init() {
	s.Server.Init()
	s.totalSize = 0
	s.Id = s.totalSize
	s.totalSize = 1
	s.newID = make(chan uint64)
}

func (s *SvrID) Start() {
	if !s.running {
		s.running = true
		go func() {
			for {
				select {
					case <-s.StopChan:
						s.running = false
						return
					case s.NewID <- s.totalSize:
						s.totalSize++
				}
			}
		}()
	}
}
