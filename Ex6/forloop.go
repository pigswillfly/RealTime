package main

import (	 
	. "fmt"
	. "strconv"
	. "time"
)

func SendMessage(chan1 chan string) {
	var buf string
	for i:=0; i<100; i++ {
		buf = (Itoa(i)+":MessageID - Timestamp: "+Time.String(Now())+" \000")
		chan1 <- buf
		Sleep(333 * Millisecond)
 	}

}

func main(){
	chan1 := make(chan string)
	var sing string
	go SendMessage(chan1)
	for {
		sing = <- chan1
		Println(sing)
		Sleep(10 * Millisecond)
	}
}
