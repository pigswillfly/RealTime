package main

import (
//	"Testing"
	."Control"
//	."fmt"
//	."strconv"
	."time"
)


func main() {
//	go Testing.Driver_Test()
//	go Testing.Network_Test()
//	go Testing.Elevator_Test()
	go Init_Control()
//	Init_Control()
	for i:=0; i<1000; i++{
		Sleep(Second)
	}
}


