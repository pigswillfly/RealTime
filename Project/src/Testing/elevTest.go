package main

import (
	"Network"
	. "fmt"
)

func main() {
	msg := "Hello World!"
	var message string
	net.Init_Net()
	net.ToNet <- msg
	message = <-net.FromNet
	Println(message)
}
