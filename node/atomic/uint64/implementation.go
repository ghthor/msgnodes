package atomic_uint64

func (an *AtomicNode) processMsg(msg Msg) {
	switch msg.(type) {
		case QueryVal:
			qv := msg.(QueryVal)
			qv.Query <- an.Val
		case QueryAndSet:
			qs := msg.(QueryAndSet)
			qs.Query <- an.Val
			an.Val = <-qs.Query
	}
}
