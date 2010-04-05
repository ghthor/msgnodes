package main

import (
	"fmt"
	"net"
	"os"
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

	raddr, err := net.ResolveTCPAddr("localhost:52000")
	if Errored(err, "While Resolving raddr") { return }

	laddr, err := net.ResolveTCPAddr("localhost:52001")
	if Errored(err, "While Resolving laddr") { return }

	socket, err := net.DialTCP("tcp", laddr, raddr)
	if Errored(err, "While attempting to connect") { return }

	bytesWritten, err := socket.Write(make([]byte, 4))
	if Errored(err, "While Writing Exit Signal to Socket") { return }
	fmt.Printf("Bytes Written: %v\n", bytesWritten)
	socket.Close()
}
