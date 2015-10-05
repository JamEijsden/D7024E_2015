package dht

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func (node *DHTNode) joinRequest(dhtNode *DHTNode) {
	rec := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//go node.sendMsg("request", rec)
	go node.join(rec)
	//retry := time.Timer(time.Millisecond.500)
	time.Sleep(time.Millisecond * 500)

	/*if node.succ[1] == "" {
		fmt.Println("NOOOPE")
		go node.sendMsg("request", rec)
	}*/

}

func TestNetwork(t *testing.T) {
	fmt.Println("Beginning test..")
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

	node0 := makeDHTNode(nil, "localhost", "1110")

	node1 := makeDHTNode(nil, "localhost", "1111")

	node2 := makeDHTNode(nil, "localhost", "1112")

	node3 := makeDHTNode(nil, "localhost", "1113")

	node4 := makeDHTNode(nil, "localhost", "1114")

	node5 := makeDHTNode(nil, "localhost", "1115")

	node6 := makeDHTNode(nil, "localhost", "1116")

	node7 := makeDHTNode(nil, "localhost", "1117")

	//	node9 := makeDHTNode(nil, "localhost", "1119")
	wg.Add(8)
	go node0.startServer(&wg)
	go node1.startServer(&wg)
	go node2.startServer(&wg)
	go node3.startServer(&wg)
	go node4.startServer(&wg)
	go node5.startServer(&wg)
	go node6.startServer(&wg)
	go node7.startServer(&wg)
	wg.Wait()

	//joinRing(node1, &wg)
	//
	node1.joinRequest(node0)

	node2.joinRequest(node0)

	node3.joinRequest(node0)
	//printFingers(node0)

	node4.joinRequest(node0)

	node5.joinRequest(node0)

	node6.joinRequest(node0)

	node7.joinRequest(node0)

	time.Sleep(time.Second * 2)
	//printFingers(node0)
	//go node4.nodeFail()

	node4.joinRequest(node0)
	time.Sleep(time.Second * 7)
	node0.QueueTask(createTask("print", createMsg("", "", "1", "", "")))
	time.Sleep(time.Second * 5)
	fmt.Println(hashString("wooot.txt"))
	node0.fingerLookup(hashString("wooot.txt"))
	//sendData(node5)
	time.Sleep(time.Second * 500)

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
