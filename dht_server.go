package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type Transport struct {
	node        *DHTNode
	bindAddress string
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	fmt.Println("Server running on " + transport.bindAddress + " with ID " + transport.node.nodeId + "\nWaiting for messages..")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err = dec.Decode(&msg)

		if msg.Type == "test" {
			fmt.Println("\n" + transport.bindAddress + " received msg from " + msg.Src)
			fmt.Println("Message type: " + msg.Type)
			fmt.Print(msg.Timestamp)
			fmt.Println(": " + msg.Key)
			reply := Msg{1337, "test", "Reply to your hello world", transport.bindAddress, msg.Src, transport.bindAddress}
			go transport.send(&reply)

		} else if msg.Type == "join" || msg.Type == "accept" {
			fmt.Println("joining..")
			transport.node.nodeJoin(&msg)
		} else {
			fmt.Println(transport.node.nodeId + ": Type operation not found for: ")
			fmt.Print(msg)

		}
		return
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
