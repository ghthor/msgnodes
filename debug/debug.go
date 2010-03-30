package debug

import (
	"fmt"
)

var chanPrintln chan string

func init() {
	chanPrintln = make(chan string, 500)
}

func PrintlnDump() {
	for !closed(chanPrintln) {
		msg := <-chanPrintln
		if msg == "" {
			msg = "I think the channel is closed"
			close(chanPrintln)
		}
		fmt.Println(msg)
	}
}

func PushMsg(msg string) {
	go func() {
		ok := chanPrintln <- msg
		if !ok {
			PrintlnDump()
			chanPrintln <- msg
		}
	}()
}
