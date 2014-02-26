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

func main(){

	// Resolve remote and local adresses
	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20007")
	laddr, err := ResolveUDPAddr("udp", ":20007")
	if err != nil{Println("Resolving address failed.")}

	// Create sockets
	conn1, err := DialUDP("udp", nil, raddr)
	if conn1 != nil{Println("UDP connection established.")}
	conn2, err := ListenUDP("udp", laddr)
	if conn2 != nil{Println("UDP listener established.")}
	
	// Closes sockets when main is done
	defer conn1.Close()
	defer conn2.Close()
    
	// Send a message
	buff := ([]byte("ASL?\000"))
	n, err := conn1.Write(buff)
	if err != nil {
	    Println("UDP write error!")
	} else {
	    Println(n, " bytes written.")
	}
	
	// Buffered channel
	c := make([]byte, 1024)
	// Listening channel fetching messages from go listen() 
	chan1 := make(chan []byte)
	// Listen after UDP messages concurrently
	go listen(chan1, conn2)
	// Wait for available message
	c = <- chan1
	
	Println("Server reply (what we sent):", string(c))
    
    //time.Sleep(10)

}

