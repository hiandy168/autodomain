package main

import "yiyecp.com/autodomain/net"

func main() {
	go net.GoRun()
	var c = make(chan int, 10)
	<-c
}
