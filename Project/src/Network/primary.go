package main;

import (	 
	. "fmt"
	. "net"
	. "strconv"
	"time"
	"os/exec"	
)

func broadcast(chan1 chan int, conn1 *UDPConn) {
	for {
		buf := ([]byte(Itoa(chan1)+"\000"))
		bytes, err := conn1.Write(buf)
		if err != nil {
		Println("UDP write error!")
		} else {
		Println("Message number: " + Itoa(chan1) + " sent, and ", bytes, " bytes written.")
		}	
		time.Sleep(333 * time.Millisecond)
	}
}

func listen(chan2 chan []byte, conn2 *UDPConn) {
	for {
		listener := make([]byte, 1024) 
		n,_,err := conn2.ReadFromUDP(listener) 
		if err!=nil{
		Println("UDP read error!")
		} 
		chan2 <-listener[:n]
	}   
}

func main(){
	raddr, err := ResolveUDPAddr("udp", "129.241.187.255:20012")
	if err != nil{Println("Resolving address failed.")}
	conn1, err := DialUDP("udp", nil, raddr)
	if conn1 != nil{Println("UDP connection established.")}
	laddr, err := ResolveUDPAddr("udp", ":20012")
	if err != nil{Println("Resolving address failed.")}
	conn2, err := ListenUDP("udp", laddr)
	if conn2 != nil{Println("UDP listener established.")}

	defer conn1.Close()	
	defer conn2.Close()

	chan1 := make(chan int)
	chan2 := make(chan []byte)
	msg := make([]byte, 1024)
	var primary bool = false
	var temp, i, j int = 0, 1, 0

	go broadcast(chan1, conn1)
	go listen(chan2, conn2)

	for {
		
		for primary == true {
			for i<1000 {
				chan1 <- i
				i++			
			}		
		}
		for primary == false {
			msg = <-chan2
			j = Atoi(string(msg))
			
			Println("Message #",j," received, at time: " + time.Time.String(time.Now()))

			if temp == j-1 || temp == j-2 || temp == j-3 {
				temp = j
			} else {
				Println("Primary process is dead!")
				Println("New primary assigned. Spawning new backup...")
				cmd := exec.Command("mate-terminal", "-x", "go", "run", "primary.go")
				cmd.Run()
				i = temp
				primary = true
				
			}
		}
		if i<=1{
			primary = true
		}
	}
}































