package dht

import (
	"fmt"
	"sync"
	"time"
	"os"
	"net"
	"strings"
)

func BootUpNode() {
	// go test -test.run TestNetork
	var wg sync.WaitGroup
	addrs, _ := net.InterfaceAddrs()
	addr := strings.Split(addrs[1].String(),"/")
//	fmt.Println(addr)
//	for _, ip := range addrs {
//		fmt.Println(ip.Network() + " " + ip.String())
//	}
	fmt.Println("Booting node with ip ", addr)
	
	node := makeDHTNode(nil, addr[0], "1110")

	wg.Add(1)
	go node.startServer(&wg)
	wg.Wait()

	node.joinRequest()
	fmt.Println("Node " + node.nodeId + " is up")

	time.Sleep(time.Second * 50000)

}

func (node *DHTNode) joinRequest() {
	bootstrap := os.Getenv("BS_PORT_8000_TCP_ADDR")
	//rec := node.contact.ip + ":" + node.contact.port
	//go node.sendMsg("request", rec)
	if "" != bootstrap {
		fmt.Println("Trying to join "+bootstrap+":1110")
		go node.join(bootstrap+":1110")
		//retry := time.Timer(time.Millisecond.500)
		time.Sleep(time.Millisecond * 100)
	}

	/*if node.succ[1] == "" {
		fmt.Println("NOOOPE")
		go node.sendMsg("request", rec)
	}*/
}

func (dhtNode *DHTNode) heartbeatRespons(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	dhtNode.transport.send(createMsg("heartbeat_respons", dhtNode.nodeId, src, msg.Src, msg.Origin))
}

func printFingers(node *DHTNode) {
	fmt.Println(node.nodeId + "> Printing fingers..")
	for _, finger := range node.fingers.fingerList {
		fmt.Println(finger)
	}
	fmt.Println("No more fingers")

}
