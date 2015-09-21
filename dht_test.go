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
		go node.helloWorld("init", rec, wg)
		wg.Wait()
		time.Sleep(time.Second * 1)
	} else {
		wg.Add(1)
		go node.helloWorld("join", rec, wg)
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
	//	node9 := makeDHTNode(nil, "localhost", "1119")
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

	node3.joinRing(node1, &wg)
	/*S
	node4.joinRing(node1, &wg)

	node5.joinRing(node1, &wg)

	node6.joinRing(node1, &wg)

	node7.joinRing(node1, &wg)

	//node2.joinRing(&wg)
	/*node3.joinRing(&wg)
	node4.joinRing(&wg)
	node5.joinRing(&wg)
	node6.joinRing(&wg)
	node7.joinRing(&wg)
	node8.joinRing(&wg)*/
	//go node2.netPrintRing(createMsg("circle", "", "localhost:1112", node2.succ[1], "localhost:1112"))

	time.Sleep(time.Second * 5)

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

/* run with go test -test.run TestDHT1 */

func TestDHT1(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0b := makeDHTNode(&id0, "localhost", "1111")
	node1b := makeDHTNode(&id1, "localhost", "1112")
	node2b := makeDHTNode(&id2, "localhost", "1113")
	node3b := makeDHTNode(&id3, "localhost", "1114")
	node4b := makeDHTNode(&id4, "localhost", "1115")
	node5b := makeDHTNode(&id5, "localhost", "1116")
	node6b := makeDHTNode(&id6, "localhost", "1117")
	node7b := makeDHTNode(&id7, "localhost", "1118")

	node0b.addToRing(node1b)
	node1b.addToRing(node2b)
	node1b.addToRing(node3b)
	node1b.addToRing(node4b)
	node4b.addToRing(node5b)
	node3b.addToRing(node6b)
	node3b.addToRing(node7b)

	fmt.Println("-> ring structure")
	node1b.printRing()

	//	fmt.Println("RESPONSIBLE: " + node1b.acceleratedLookupUsingFingerTable("04", 2).nodeId)
	//findFingers(node0b)
	//node3b.testCalcFingers(0, 3) //node 4, dist 1
	//node3b.printFingers(1, 3) // node 4, dist 1
	//node3b.printFingers(2, 3) // node 5, dist 2
	//node3b.printFingers(3, 3) // node 7, dist 4

}

func TestBetween(t *testing.T) {
	testBetween("03", "02", "02")
}
func TestDistance(t *testing.T) {
	fmt.Println(distance([]byte("03"), []byte("06"), 3))
}
func TestDHT2(t *testing.T) {
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	fmt.Println("-> ring structure")
	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)

	node1.printRing()

	key1 := "2b230fe12d1c9c60a8e489d028417ac89de57635"
	key2 := "87adb987ebbd55db2c5309fd4b23203450ab0083"
	key3 := "74475501523a71c34f945ae4e87d571c2c57f6f3"

	fmt.Println("TEST: " + node1.lookup(key1).nodeId + " is responsible for " + key1)
	fmt.Println("TEST: " + node1.lookup(key2).nodeId + " is responsible for " + key2)
	fmt.Println("TEST: " + node1.lookup(key3).nodeId + " is responsible for " + key3)

	nodeForKey1 := node1.acceleratedLookupUsingFingerTable(key1)
	fmt.Println("dht node " + nodeForKey1.nodeId + " running at " + nodeForKey1.contact.ip + ":" + nodeForKey1.contact.port + " is responsible for " + key1)
	nodeForKey2 := node1.acceleratedLookupUsingFingerTable(key2)
	fmt.Println("dht node " + nodeForKey2.nodeId + " running at " + nodeForKey2.contact.ip + ":" + nodeForKey2.contact.port + " is responsible for " + key2)
	nodeForKey3 := node1.acceleratedLookupUsingFingerTable(key3)
	fmt.Println("dht node " + nodeForKey3.nodeId + " running at " + nodeForKey3.contact.ip + ":" + nodeForKey3.contact.port + " is responsible for " + key3)
}
