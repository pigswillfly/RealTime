package main;

import (	 
	. "fmt"
	. "net"
	"time"
	. "strconv"
	"os/exec"
)

func SendMessage() {

}

func main(){
	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20019")
	if err != nil{Println("Resolving address failed.")}
	conn1, err := DialUDP("udp", nil, raddr)
	if conn1 != nil{Println("UDP connection established.")}
	defer conn1.Close()
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "backup1.go")
	cmd.Run()
	
	// Sends 3 messages/second. The message has a number ID and a timestamp. 
	for i:=0; i<100; i++ {
		// Sends message on the format 1:MessageID - Timestamp: yyyy-mm-dd hh:mm:ss.ns +0100 CET
		buf := ([]byte(Itoa(i)+":MessageID - Timestamp: "+time.Time.String(time.Now())+" \000"))
		bytes, err := conn1.Write(buf)
		if err != nil {
		    Println("UDP write error!")
		} else {
		    Println("Message number: " + Itoa(i) + " sent, and ", bytes, " bytes written.")
		}
		
		time.Sleep(333 * time.Millisecond)
 	}

}
