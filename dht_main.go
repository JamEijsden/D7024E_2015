package dht

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// go test -test.run TestNetork
	fmt.Println("Booting node..")
	var wg sync.WaitGroup
	/*
		id1 := "00"
		id2 := "01"
		id3 := "02"
		id4 := "03"
		id5 := "04"
		id6 := "05"
		id7 := "06"
		id8 := "07"
			node0 := makeDHTNode(&id1, "localhost", "1110")

			node1 := makeDHTNode(&id2, "localhost", "1111")

			node2 := makeDHTNode(&id3, "localhost", "1112")

			node3 := makeDHTNode(&id4, "localhost", "1113")

			node4 := makeDHTNode(&id5, "localhost", "1114")

			node5 := makeDHTNode(&id6, "localhost", "1115")

			node6 := makeDHTNode(&id7, "localhost", "1116")

			node7 := makeDHTNode(&id8, "localhost", "1117")

	*/
	node := makeDHTNode(nil, "localhost", "1110")

	wg.Add(1)
	go node.startServer(&wg)
	wg.Wait()

	node.joinRequest()
	fmt.Println("Node " + node.nodeId + " is up")

	time.Sleep(time.Second * 5000)

}

func (node *DHTNode) joinRequest() {

	rec := node.contact.ip + ":" + node.contact.port
	//go node.sendMsg("request", rec)
	if "localhost:1110" != rec {
		go node.join(rec)
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
