package main

import (
	"fmt"
	//"net"
	"os"
	"time"
	"runtime"
	"rand"
	//"container/vector"
	//"container/list"
	//"rand"
	//i "ghthor/init"
	//"ghthor/node"
	//bufNode "ghthor/node/buffer"
	//"ghthor/node"
)

func Errored(e os.Error, mesg string) bool {
	if(e != nil) {
		fmt.Printf("Errored %s\n",mesg)
		fmt.Println(e.String())
		return true
	}
	return false
}

func main() {
	runtime.GOMAXPROCS(4)
	ts, _, _ := os.Time()
	fmt.Printf("TimeSec: %v\n", ts)
	rand.Seed(ts)
	fmt.Printf("\n\n");


	joinCh := make(chan bool)
	defer close(joinCh)

	var Uint uint64 = 100
	off := -1
	Uint += off
	//<-joinCh
	time.Sleep(2000)
	fmt.Printf("\n\n")
}
