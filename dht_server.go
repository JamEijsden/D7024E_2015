package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type Transport struct {
	node        *DHTNode
	bindAddress string
	queue chan *Msg
}

func CreateTransport(node *DHTNode, bindAddress string) {
	transport := &Transport{}
	transport.queue make(chan *Msg)
	transport.node = node
	transport.bindAddress = bindAddress

	return transport
}

func (transport *Transport) processMsg() {
	msg := <- transport.queue

	go func() {
		  for {
      		select {
      			case msg1 := <- c1:
        fmt.Println(msg1)
      case msg2 := <- c2:
        fmt.Println(msg2)
      }
    }
	}()
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)  //adds adress to variable udpAddr and err msg in var err is there is one.
	conn, err := net.ListenUDP("udp", udpAddr) //we listen to the IP-adress in udpAddr.
	fmt.Println("Server running on " + transport.bindAddress + " with ID " + transport.node.nodeId + "\nWaiting for messages..")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	dec := json.NewDecoder(conn) //decodes conn and adds to variable dec
	for {
		msg := Msg{}
		err = dec.Decode(&msg) //decodes the message where the message adress is and adds an error msg if there's none.

		transport.queue <- msg

	
		//return
		//	we	got	a	message, do something
	}
}

func (transport *Transport) send(msg *Msg) {
	fmt.Println("\nPreparing to send msg from " + msg.Src + " to " + msg.Dst)
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(msg)
	bytes, err := json.Marshal(msg)
	defer conn.Close()
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println(err.Error())
	}

}
