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
	Msgs []list.List
	MsgsBuffered uint64
	MsgsProcessed uint64
	BufferTo uint64
}

func (bn *BufferNode) Init(bufferSz int, ShutDownCh chan int) (*BufferNode) {
	bn.BaseNode.Init(ShutDownCh)
	bn.In = make(chan Msg, bufferSz + 1)
	bn.MsgReq = make(chan (chan Msg))
	bn.Msgs = make([]list.List, 10)
	bn.MsgsBuffered = 0
	bn.MsgsProcessed = 0
	bn.BufferTo = uint64(bufferSz)
	return bn
}

func (bn *BufferNode) Dispose() {
	close(bn.MsgReq)
	close(bn.In)
	bn.BaseNode.Dispose()
}

func(bn *BufferNode) Listen() {
	bn.Lock()
	if !bn.Running {
		bn.Running = true
		bn.Unlock()
		go func() {
			MsgChan := make(chan Msg)
			for !closed(bn.In) || !closed(bn.ShutDownCh) {
				select {
					case sdVal := <-bn.ShutDownCh:
						bn.ShutDown(sdVal)
						return
					//case msg := <-bn.In:
						//bn.bufferMsg(msg)
					case bn.MsgReq <- MsgChan:
						// in theory I think this can be written like this
						//if bn.Fill() && bn.FillTo(1) {
						//	bn.PopNextInto(MsgChan)
						//} else {
						//	MsgChan <- nil // An Error
						//}
						bn.PopMsgInto(MsgChan)
				}
			}
		}()
	}
}

func (bn *BufferNode) PopMsgInto(MsgChan chan Msg) {
	if bn.Fill() {
		if bn.MsgsBuffered > 0 {
			bn.PopNextInto(MsgChan)
		} else if bn.FillTo(1) {
			// I think this is just being paranoid
			if bn.MsgsBuffered > 0 {
				bn.PopNextInto(MsgChan)
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

func (bn *BufferNode) PopNextInto(MsgChan chan Msg) {
	// i represents the Priority, 0 being the Highest
	for i := 0; i < len(bn.Msgs); i++ {
		// If there is a Msg at this Priority Level
		if bn.Msgs[i].Len() > 0 {
			ele := bn.Msgs[i].Front()
			MsgChan <- ele.Value.(Msg)
			bn.Msgs[i].Remove(bn.Msgs[i].Front())
			bn.MsgsBuffered--
			break
		}
	}
}

func (bn *BufferNode) FillTo(numToBuf uint64) (bool) {
	for bn.MsgsBuffered < numToBuf {
		select {
			case msg := <-bn.In:
				bn.bufferMsg(msg)
			case sdVal := <-bn.ShutDownCh:
				bn.ShutDown(sdVal)
				return false
		}
	}
	// Some Sort of Error, I dunno
	if bn.MsgsBuffered < numToBuf {
		return false
	}
	return true
}

func (bn *BufferNode) Fill() (bool) {
	for bn.MsgsBuffered < bn.BufferTo {
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
	msg.SetRecvId(bn.MsgsProcessed)
	bn.MsgsProcessed++
	bn.MsgsBuffered++
	bn.Msgs[msg.Priority()].PushBack(msg)
}
