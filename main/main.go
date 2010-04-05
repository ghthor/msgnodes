package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"rand"
	//"container/vector"
	//"container/list"
	//"rand"
	//i "ghthor/init"
	//"ghthor/node"
	bufNode "ghthor/node/buffer"
	"ghthor/node"
)

func Errored(e os.Error, mesg string) bool {
	if(e != nil) {
		fmt.Printf("Errored %s\n",mesg)
		fmt.Println(e.String())
		return true
	}
	return false
}

type HighPrior struct {
	node.BaseMsg
}

func (hp *HighPrior) Priority() (uint64) {
	return 0
}

type LowPrior struct {
	node.BaseMsg
}

func (lp *LowPrior) Priority() (uint64) {
	return 2
}

type Msg node.Msg

func main() {
	runtime.GOMAXPROCS(4)
	ts, _, _ := os.Time()
	fmt.Printf("TimeSec: %v\n", ts)
	rand.Seed(ts)
	fmt.Printf("\n\n");

	buf := new(bufNode.BufferNode).Init(40, nil)
	defer buf.Dispose()
	buf.Listen()

	newMsg := make(chan Msg, 2)
	defer close(newMsg)

	// Feed the Buffer
	go func() {
		for !closed(newMsg) {
			switch rand.Intn(10) {
				case 0:
					//newMsg <- HighPrior{}.(Msg)
					newMsg <- new(HighPrior)
				case 1,2,3,4,5,6,7,8,9:
					//newMsg <- LowPrior{}.(Msg)
					newMsg <- new(LowPrior)
			}
		}
	}()

	go func() {
		// Fill Buffer with 10 Msgs
		for i := 0; i < 200; i++ {
			if closed(newMsg) || closed(buf.In) { return }
			msg := <-newMsg
			buf.In <- msg
		}
	}()

	var MsgsProcessed uint64 = 0

	go func() {
		for !closed(buf.MsgReq) {
			msgChan := <-buf.MsgReq
			if msgChan == nil { return }
			msg := <-msgChan
			if msg == nil { return }
			msg.SetProcId(MsgsProcessed)
			MsgsProcessed++
			fmt.Printf("MsgPrior: %v \tRecvId: %v \tProcId: %v \tMsgsBuffered: %v\n", msg.Priority(), msg.RecvId(), msg.ProcId(), buf.MsgsBuffered)
		}
	}()

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
