package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"container/vector"
	//"rand"
	//i "ghthor/init"
	//"ghthor/node"
)

func NewAtomInt64Node(in chan interface {}) (*AtomInt64Node) {
	an := new(AtomInt64Node)
	an.shutDown = make(chan int, 1)
	an.isRunning = false
	if in == nil {
		an.in = make(chan interface {})
	} else {
		an.in = in
	}
	an.Listen()
	return an
}

type AtomInt64Node struct {
	val int64
	isRunning bool
	in chan interface {}
	shutDown chan int
	monitors vector.Vector
}

type Query struct {
	val chan int64
}

type SetTo struct {
	val int64
}

type OffsetBy struct {
	val int64
}

type QueryAndSet struct {
	Query
}

type AddMonitor struct {
	monitor chan interface {}
}

type RemoveMonitor struct {
	monitor *AddMonitor
}

func(an *AtomInt64Node) Stop() {
	go func() {
		if !closed(an.shutDown) {
			an.shutDown <- 0
			close(an.shutDown)
		}
	}()
}

func (an *AtomInt64Node) Listen() {
	if !an.isRunning {
		an.isRunning = true
		go func() {
			for !closed(an.in) && !closed(an.shutDown) {
				select {
					case msg := <-an.in:
						an.processMsg(msg)
					case sdVal := <-an.shutDown:
						an.isRunning = false
						sdVal++
						if !closed(an.shutDown) { an.shutDown <- sdVal }
						return
				}
			}
		}()
	}
}

func (an *AtomInt64Node) processMsg(msg interface {}) {
	switch msg.(type) {
		case nil:
			//fmt.Println("Nil Msg")
		case Query:
			query := msg.(Query)
			query.val <- an.val
		case SetTo:
			setTo := msg.(SetTo)
			an.val = setTo.val
			an.informMonitors(an.val)
		case QueryAndSet:
			qAndSet := msg.(QueryAndSet)
			qAndSet.val <- an.val
			an.val = <-qAndSet.val
			an.informMonitors(an.val)
		case OffsetBy:
			offsetBy := msg.(OffsetBy)
			an.val += offsetBy.val
			an.informMonitors(an.val)
		case *AddMonitor:
			//an.monitors.Push(msg.(*AddMonitor))
			an.monitors.Push(msg)
		case RemoveMonitor:
			removee := msg.(RemoveMonitor).monitor
			for i, ele := range an.monitors.Data() {
				//elem, ok := ele.(*AddMonitor)
				elem := ele.(*AddMonitor)
				if elem == removee {
					an.monitors.Delete(i)
					return
				}
			}
	}
}

func (an *AtomInt64Node) informMonitors(msg interface {}) {
	for _, ele := range an.monitors.Data() {
		elem := ele.(*AddMonitor)
		elem.monitor <- msg
	}
}

func Errored(e os.Error, mesg string) bool {
	if(e != nil) {
		fmt.Printf("Errored %s\n",mesg)
		fmt.Println(e.String())
		return true
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(2)
	fmt.Printf("\n\n");

	statusMon := make(chan interface {}, 500)
	defer close(statusMon)

	counter := NewAtomInt64Node(make(chan interface {}))
	defer counter.Stop()

	go func() {
		queryChan := make(chan int64)
		for !closed(statusMon) {
			newVal := <-statusMon
			if newVal == nil { return }
			fmt.Printf("Counter is Now: %v\n", newVal.(int64))
			if newVal.(int64) == -5 || newVal.(int64) == 5 {
				fmt.Println("Setting to 0")
				counter.in <- QueryAndSet{Query{val:queryChan}}
				curVal := <-queryChan
				if curVal == newVal.(int64) {
					queryChan <- 0
				} else {
					queryChan <- 0
				}
				fmt.Println("Counter Should be 0 Now")
			}
		}
	}()

	randChan := make(chan int64)
	defer close(randChan)

	go func() {
		for !closed(randChan) {
			select {
				case randChan <- 1:
				case randChan <- -1:
			}
		}
	}()

	go func() {
		for i := 0; i < 300; i++ {
			val := <-randChan
			counter.in <- OffsetBy{val:val}
		}
	}()

	counter.in <- &AddMonitor{monitor:statusMon}
	counter.in <- SetTo{val:-10}

	listenAddr, err := net.ResolveTCPAddr("localhost:52000")
	if Errored(err, "Resolving TCP Addr") { return }

	fmt.Println("Listening For Connections")
	tcpServer, err := net.ListenTCP("tcp4",listenAddr)
	if Errored(err, "Starting to Listen for Connections") { return }

	fmt.Println("Waiting on Connection")
	socket, err := tcpServer.AcceptTCP()
	if Errored(err, "While Accepting TCP Connection") { return }
	tcpServer.Close()

	data := make([]byte, 4)
	fmt.Println("Waiting For Kill Signal")
	bytesRead, err := socket.Read(data)
	fmt.Printf("BytesRead: %v", bytesRead)
	if Errored(err, "While Waiting for the kill Signal") { return }
	socket.Close()

	fmt.Printf("\n\n");
}
