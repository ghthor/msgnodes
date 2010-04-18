package buffer

import (
	"ghthor/node"
	"container/list"
)

type Msg node.Msg

type BufferNode struct {
	node.BaseNode
	In chan Msg
	MsgReq chan (chan Msg)
	msgs []list.List
	msgsBuffered uint64
	msgsProcessed uint64
	bufferTo uint64
}

func (bn *BufferNode) Init(bufferSz int, ShutDownCh chan int) (*BufferNode) {
	bn.BaseNode.Init(ShutDownCh)
	bn.In = make(chan Msg, bufferSz + 1)
	bn.MsgReq = make(chan (chan Msg))
	bn.msgs = make([]list.List, 10)
	bn.msgsBuffered = 0
	bn.msgsProcessed = 0
	bn.bufferTo = uint64(bufferSz)
	return bn
}

func (bn *BufferNode) Dispose() {
	close(bn.MsgReq)
	close(bn.In)
	bn.BaseNode.Dispose()
}

func(bn *BufferNode) IsRunning() bool {
	bn.Lock()
	defer bn.Unlock()
	return bn.Running
}

func(bn *BufferNode) Listen() {
	bn.Lock()
	if !bn.Running {
		bn.Running = true
		bn.Unlock()
		go func() {
			MsgChan := make(chan Msg)
			defer close(MsgChan)
			for !closed(bn.In) || !closed(bn.ShutDownCh) {
				select {
					case sdVal := <-bn.ShutDownCh:
						bn.ShutDown(sdVal)
						return
					//case msg := <-bn.In:
						//bn.bufferMsg(msg)
					case bn.MsgReq <- MsgChan:
						// in theory I think this can be written like this
						//if bn.fill() && bn.fillTo(1) {
						//	bn.PopNextInto(MsgChan)
						//} else {
						//	MsgChan <- nil // An Error
						//}
						bn.popMsgInto(MsgChan)
				}
			}
		}()
	}
}

func (bn *BufferNode) popMsgInto(MsgChan chan Msg) {
	if bn.fill() {
		if bn.msgsBuffered > 0 {
			bn.popNextInto(MsgChan)
		} else if bn.fillTo(1) {
			// I think this is just being paranoid
			if bn.msgsBuffered > 0 {
				bn.popNextInto(MsgChan)
			} else { // Some Weird Error or ShutDown
				MsgChan <- nil
			}
		} else { // Some Weird Error or ShutDown
			MsgChan <- nil
		}
	} else { // Some Weird Error or ShutDown
		MsgChan <- nil
	}
}

func (bn *BufferNode) popNextInto(MsgChan chan Msg) {
	// i represents the Priority, 0 being the Highest
	for i := 0; i < len(bn.msgs); i++ {
		// If there is a Msg at this Priority Level
		if bn.msgs[i].Len() > 0 {
			ele := bn.msgs[i].Front()
			MsgChan <- ele.Value.(Msg)
			bn.msgs[i].Remove(bn.msgs[i].Front())
			bn.msgsBuffered--
			break
		}
	}
}

func (bn *BufferNode) fillTo(numToBuf uint64) (bool) {
	for bn.msgsBuffered < numToBuf {
		if closed(bn.In) || closed(bn.ShutDownCh) { return false }
		select {
			case msg := <-bn.In:
				bn.bufferMsg(msg)
			case sdVal := <-bn.ShutDownCh:
				bn.ShutDown(sdVal)
				return false
		}
	}
	// Some Sort of Error, I dunno
	if bn.msgsBuffered < numToBuf {
		return false
	}
	return true
}

func (bn *BufferNode) fill() (bool) {
	for bn.msgsBuffered < bn.bufferTo {
		if closed(bn.In) || closed(bn.ShutDownCh) { return false }
		select {
			case msg := <-bn.In:
				bn.bufferMsg(msg)
			case sdVal := <-bn.ShutDownCh:
				bn.ShutDown(sdVal)
				return false
			default:
				goto Success
		}
	}
	Success:
	return true
}

// Buffer the Msg into the Array of Lists
func (bn *BufferNode) bufferMsg(msg Msg) {
	msg.SetRecvId(bn.msgsProcessed)
	bn.msgsProcessed++
	bn.msgsBuffered++
	bn.msgs[msg.Priority()].PushBack(msg)
}
