
package main;

import (
	. "fmt"
	. "runtime"
	. "time"
)

var i = 0;

func adder(){
	for x:=0; x<1000000; x++{
		i++
	}
}

func subber(){
	for x:=0; x<1000000; x++{
		i--
	}
}

func main(){
	GOMAXPROCS(NumCPU())
	go adder()
	go subber()

	Sleep(100*Millisecond)
	Println("Done:", i)
}
