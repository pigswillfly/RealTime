
package main;

import (
	. "fmt"
	. "runtime"
	. "time"
)

var i = 0
var ci = make(chan int)

func adder(){
	for x:=0; x<1000000; x++{
		<- ci
		i++
		ci <- 1
	}
}

func subber(){
	for x:=0; x<1000000; x++{
		<- ci
		i--
		ci <- 1
	}
}

func main(){
	GOMAXPROCS(NumCPU())
	
	go adder()
	go subber()

	Sleep(100*Millisecond)
	Println("Done:", i)
}
