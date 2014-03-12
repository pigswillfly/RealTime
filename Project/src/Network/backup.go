package main

import (	 
	. "fmt"
	. "net"
	 //"time"
)



/*func backup(){

}*/

func main(){
	laddr, err := ResolveUDPAddr("udp", ":20012")
	if err != nil{Println("Resolving address failed.")}
	sock, err := ListenUDP("udp", laddr)
	if sock != nil{Println("UDP listener established.")}
	defer conn2.Close()
	chan1 := make(chan string)
	var msg string
	go listen(sock, chan1)

	for {
		msg = <- chan1
		Println(string(msg))
	}

}
