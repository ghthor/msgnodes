package main

import (
	"fmt"
	"os"
	"time"
	"runtime"
	"rand"
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
	defer buf.Stop()
	//defer buf.Dispose()
	buf.Listen()

	newMsg := make(chan Msg, 2)
	defer close(newMsg)

	// Feed the Buffer
	go func() {
		for !closed(newMsg) {
			switch rand.Intn(10) {
				case 0:
					newMsg <- new(HighPrior)
				case 1,2,3,4,5,6,7,8,9:
					newMsg <- new(LowPrior)
			}
		}
	}()

	joinCh := make(chan bool)
	defer close(joinCh)

	numMsgs := 200
	go func() {
		// Fill Buffer with 10 Msgs
		for i := 0; i < numMsgs; i++ {
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
			if MsgsProcessed >= uint64(numMsgs) { joinCh <- true; return }
		}
	}()


	<-joinCh
	time.Sleep(2000)
	fmt.Printf("\n\n")
}
