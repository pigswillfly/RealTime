package Network

import (
	. "fmt"
	. "net"
	"time"
	"os"
)

const (
	myPort = "24642"
	subNet = "129.241.187.255"
	chanSize = 1024
)

type Network struct {
	FromNet chan string
	ToNet chan string 
}

func Init_Net(to chan string, from chan string) *Network {
	net := new(Network)
	net.ToNet = to
	net.FromNet = from
	go net.SendUDP()
	go net.ReadUDP()

	return net
}

func (net *Network) SendUDP() {
	buf := make([]byte, 1024)

	raddr, err := ResolveUDPAddr("udp", subNet+":"+myPort)
	if err != nil {
		Println("Resolving address failed.")
		os.Exit(-1)
	} 
	rconn, err := DialUDP("udp", nil, raddr)
	if err != nil {
		Println("UDP dialup failed.")
		os.Exit(-2)
	} 

	for {
		msg := <-net.ToNet
	//	Println(msg)
		buf = []byte(msg)
		rconn.Write(buf)
		time.Sleep(time.Millisecond * 5) // Change frequency here.
	}

	rconn.Close()
}

func (net *Network) ReadUDP() {
	buf := make([]byte, 1024)
	var n int 

	laddr, err := ResolveUDPAddr("udp", ":"+myPort)
	if err != nil {
		Println("Resolving address failed.")
		os.Exit(-3)
	}

	lconn, err := ListenUDP("udp", laddr)
	if err != nil {
		Println("UDP listener failed to establish.")
		os.Exit(-4)
	} 

	for {
		 n, _ = lconn.Read(buf)
		// Println(string(buf[:n]))
		 net.FromNet <- string(buf[:n])
		 time.Sleep(time.Millisecond * 5)
	}

}


