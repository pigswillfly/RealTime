package Network

import (
	. "fmt"
	. "net"
	"time"
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

func Init_Net() *Network {
	net := new(Network)
	net.FromNet = make(chan string, chanSize)
	net.ToNet = make(chan string, chanSize)
	go net.SendUDP()
	go net.ReadUDP()

	return net
}

func (net *Network) SendUDP() {
	buf := make([]byte, 1024)

	raddr, err := ResolveUDPAddr("udp", subNet+":"+myPort)
	if err != nil {
		Println("Resolving address failed.")
	}
	rconn, err := DialUDP("udp", nil, raddr)
	if err != nil {
		Println("UDP dialup failed.")
	}

	for {
		msg := <-net.ToNet
		buf = []byte(msg)
		rconn.Write([]byte(buf))
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
	}

	lconn, err := ListenUDP("udp", laddr)
	if err != nil {
		Println("UDP listener failed to establish.")
	}

	for {
		 n, _ = lconn.Read(buf)
		 net.FromNet <- string(buf[:n])
		 time.Sleep(time.Millisecond * 5)
	}

}

/*for {
	select{
	case newOrder := <-newOrderChan:
		//Handle new order
	case someoneRequestsSomethingChan<- something:
	case <-closeChan:
		return
	default:
		// do something
	}
	
}*/
