package main;

import (	 
	. "fmt"
	. "net"
	"time"
)

func tcp(){
// Set message to size 1024, if port 34933
// \0 as the marker, if you connect to port 33546
// Connect to: 129.241.187.157\0
	
	// Remote address
	raddr, err := ResolveTCPAddr("tcp", "129.241.187.161:33546")
	if err != nil{}

	// Dial server
	conn, err := DialTCP("tcp", nil, raddr)

	// Read server response
	buff := make([]byte, 1024)
	conn.Read(buff)
	s := string(buff[0:1023])
	Println("Server says :", s)
	
	// Local address
	laddr, err := ResolveTCPAddr("tcp", "129.241.187.157:22222")

	// Create listener
	listen, err := ListenTCP("tcp", laddr)

	// Ask server to connect
	buff2 := ([]byte("Connect to: 129.241.187.157:22222\000"))
	_, err = conn.Write(buff2)

	// Accept connection
	conn2, err := listen.AcceptTCP()

	// Send a message to server
	buff2 = ([]byte("Hi there!\000"))
	_, err = conn2.Write(buff2)
	
	// Read server response
	conn2.Read(buff[0:1023])
	s = string(buff[0:1023])
	Println(s)

	time.Sleep(1)

}

func udp(conn *UDPConn){
	
	buff := ([]byte("Sweet sweet buffer\000"))	
	
	for{
		n, err := conn.Read(buff)
		if err != nil{}
		Println(n, " bytes received")
		s := string(buff)
		Println(s)
	}

}


func main(){

	// Remote address
	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20019")
	if err != nil{}

	// Local address
	laddr, err := ResolveUDPAddr("udp", "129.241.187.157:0")

	// Dial server
	conn, err := DialUDP("udp", laddr, raddr)
	if conn != nil{
		Println("UDP connection established")
	}

	// Create listener
	listen, err := ListenUDP("udp", laddr)
	if listen != nil {
		Println("UDP listener created")
	}

	// Send a message to server
	buff := ([]byte("Hi there!\000"))
	n, err := conn.Write(buff)
	Println(n, " bytes written")

	// Read server response
	n, err = conn.Read(buff)
	Println(n, " bytes received")
	s := string(buff)
	Println(s)
	
	time.Sleep(10)


}



