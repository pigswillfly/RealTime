package main

import (	 
	. "fmt"
	. "net"
	 //"time"
)

func listen(chan1 chan []byte, conn2 *UDPConn){
    for { 
        listener := make([]byte, 1024) 
        n,_,err := conn2.ReadFromUDP(listener) 
        if err!=nil{ //handle error 
        } 
        chan1 <-listener[:n]     
    }

}

/*func backup(){

}*/

func main(){
	laddr, err := ResolveUDPAddr("udp", ":20019")
	if err != nil{Println("Resolving address failed.")}
	conn2, err := ListenUDP("udp", laddr)
	if conn2 != nil{Println("UDP listener established.")}
	defer conn2.Close()
	chan1 := make(chan []byte)
	go listen(chan1, conn2)
	// go backup(id, timestamp, channel)

	for {
		msg := <- chan1
		Println(string(msg))
	}

}
