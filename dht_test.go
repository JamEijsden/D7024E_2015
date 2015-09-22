package dht

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func (node *DHTNode) joinRing(dhtNode *DHTNode, wg *sync.WaitGroup) {
	rec := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.succ[0] == "" {
		//Initiate ring
		wg.Add(1)
		go node.sendMsg("init", rec, wg)
		wg.Wait()
		time.Sleep(time.Second * 1)
		if node.succ[0] == "" {
			node.joinRing(dhtNode, wg)
		}
	} else {
		wg.Add(1)
		go node.sendMsg("join", rec, wg)
		wg.Wait()
	}
}

func initiateRing() {

}
func TestNetwork(t *testing.T) {
	fmt.Println("Beginning test..")
	var wg sync.WaitGroup

	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"
	id8 := "08"

	node1 := makeDHTNode(&id1, "localhost", "1111")
	node2 := makeDHTNode(&id2, "localhost", "1112")
	node3 := makeDHTNode(&id3, "localhost", "1113")
	node4 := makeDHTNode(&id4, "localhost", "1114")
	node5 := makeDHTNode(&id5, "localhost", "1115")
	node6 := makeDHTNode(&id6, "localhost", "1116")
	node7 := makeDHTNode(&id7, "localhost", "1117")
	node8 := makeDHTNode(&id8, "localhost", "1118")
	/*
		node1 := makeDHTNode(nil, "localhost", "1111")
		node2 := makeDHTNode(nil, "localhost", "1112")
		node3 := makeDHTNode(nil, "localhost", "1113")
		node4 := makeDHTNode(nil, "localhost", "1114")
		node5 := makeDHTNode(nil, "localhost", "1115")
		node6 := makeDHTNode(nil, "localhost", "1116")
		node7 := makeDHTNode(nil, "localhost", "1117")
		node8 := makeDHTNode(nil, "localhost", "1118")
	*/ //	node9 := makeDHTNode(nil, "localhost", "1119")
	wg.Add(7)
	go node1.startServer(&wg)
	go node2.startServer(&wg)
	go node3.startServer(&wg)
	go node4.startServer(&wg)
	go node5.startServer(&wg)
	go node6.startServer(&wg)
	go node7.startServer(&wg)
	wg.Wait()

	//joinRing(node1, &wg)
	node2.joinRing(node1, &wg)

	//

	//	node3.joinRing(node1, &wg)

	node4.joinRing(node1, &wg)

	//	node5.joinRing(node1, &wg)

	node6.joinRing(node1, &wg)

	node7.joinRing(node1, &wg)

	time.Sleep(time.Second * 1)
	findFingers(node1)
	node1.fingerLookup("06")

	fmt.Println(node1.contact.port + "> " + node1.pred[0] + " " + node1.succ[0])
	fmt.Println(node2.contact.port + "> " + node2.pred[0] + " " + node2.succ[0])
	fmt.Println(node3.contact.port + "> " + node3.pred[0] + " " + node3.succ[0])
	fmt.Println(node4.contact.port + "> " + node4.pred[0] + " " + node4.succ[0])
	fmt.Println(node5.contact.port + "> " + node5.pred[0] + " " + node5.succ[0])
	fmt.Println(node6.contact.port + "> " + node6.pred[0] + " " + node6.succ[0])
	fmt.Println(node7.contact.port + "> " + node7.pred[0] + " " + node7.succ[0])

	//node1.netPrintRing(nil)

	node8.transport.listen()

	/*fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	//fmt.Print(input)

	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	if string(b) == "q" {
		return
	} */
}
