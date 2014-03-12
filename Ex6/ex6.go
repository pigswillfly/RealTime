package main;

import (	 
	. "fmt"
	. "net"
	. "strconv"
	"time"
	"os/exec"	
)

func main(){
	// Backup:
	laddr, err := ResolveUDPAddr("udp", ":20012")
	if err != nil{Println("Resolving address failed.")}

	udprecv, err := ListenUDP("udp", laddr)
	if udprecv != nil{Println("UDP listener established.")}

	buf := make([]byte, 1024) 
	var i int = 0

	for {
		udprecv.SetReadDeadline(time.Now().Add(1000*time.Millisecond))
		//read
		n,_,err := udprecv.ReadFromUDP(buf) 
		if err != nil {
			udprecv.Close()
			Println("Receive from UDP fail.")
			break
		} else {
			i, _ = Atoi(string(buf[0:n]))
		}
	}



	// Spawn duplicate:
	cmd := exec.Command("mate-terminal", "-x", "go", "run", "ex6.go")
	cmd.Run()

	// Master:
	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20012")
	if err != nil{Println("Resolving address failed.")}

	udpsend, err := DialUDP("udp", nil, raddr)
	if udpsend != nil{Println("UDP connection established.")}

	for {
		i++
		Println("Message number#",i, " at time:" + time.Time.String(time.Now()))
		buf = ([]byte(Itoa(i)))
		_, err := udpsend.Write(buf)
		if err != nil {
			Println("UDP write error!")
		} 
		time.Sleep(333 * time.Millisecond)
	}
}


