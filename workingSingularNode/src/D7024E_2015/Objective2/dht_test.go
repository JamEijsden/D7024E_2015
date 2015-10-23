package dht

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func (node *DHTNode) joinRequest(dhtNode *DHTNode) {
	rec := dhtNode.contact.ip + ":" + dhtNode.contact.port
	go node.sendMsg("request", rec)
	time.Sleep(time.Millisecond * 10)
	/*if node.succ[1] == "" {
		fmt.Println("NOOOPE")
		go node.sendMsg("request", rec)
	}*/

}

func TestNetwork(t *testing.T) {
	fmt.Println("Beginning test..")
	var wg sync.WaitGroup

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
	node1.joinRequest(node0)

	//

	node2.joinRequest(node0)

	node3.joinRequest(node0)

	node5.joinRequest(node0)

	node4.joinRequest(node0)
	//time.Sleep(time.Millisecond * 3000)
	node6.joinRequest(node0)
	node7.joinRequest(node0)
	//node0.QueueTask(createTask("print", createMsg("", "", "1", "", "")))

	time.Sleep(time.Second * 10)
	node5.leaveRing()

	/*fmt.Println(node1)
	fmt.Println(node2)
	fmt.Println(node3)
	fmt.Println(node4)
	fmt.Println(node5)
	fmt.Println(node6)
	fmt.Println(node7)
	*/
	//findFingers(node1)
	//printFingers(findFingers(node2))
	//printFingers(node2.fingers)
	//node1.fingerLookup("06")
	/*
		fmt.Println(node1.contact.port + "> " + node1.pred[0] + " " + node1.succ[0])
		fmt.Println(node2.contact.port + "> " + node2.pred[0] + " " + node2.succ[0])
		fmt.Println(node3.contact.port + "> " + node3.pred[0] + " " + node3.succ[0])
		fmt.Println(node4.contact.port + "> " + node4.pred[0] + " " + node4.succ[0])
		fmt.Println(node5.contact.port + "> " + node5.pred[0] + " " + node5.succ[0])
		fmt.Println(node6.contact.port + "> " + node6.pred[0] + " " + node6.succ[0])
		fmt.Println(node7.contact.port + "> " + node7.pred[0] + " " + node7.succ[0])
		//updateFingers(node1)
	*/
	//node1.netPrintRing(nil)

	time.Sleep(time.Second * 500)
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

func printFingers(node *DHTNode) {
	fmt.Println(node.nodeId + "> Printing fingers..")
	for _, finger := range node.fingers.fingerList {
		fmt.Println(finger)
	}
	fmt.Println("No more fingers")

}
