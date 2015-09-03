package dht

import (
	"encoding/hex"
	"fmt"
)

const BITS int = 3

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	fingers     *FingerTable
	contact     Contact
}

type FingerTable struct {
	start          string
	fingerTable    [BITS]*DHTNode
	fingerDistance [BITS]int64
	bits           int
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	dhtNode.fingers = new(FingerTable)
	dhtNode.fingers.bits = BITS
	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}
	//dhtNode.fingers.fingerTable = createFingers(dhtNode)
	dhtNode.successor = nil
	dhtNode.predecessor = nil

	return dhtNode
}

func testBetween(id1, id2, key string) {
	fmt.Println(id1 + " " + id2 + " " + key + " ")
	fmt.Println(between([]byte(id1), []byte(id2), []byte(key)))
}

/* Connects and rearranges nodes */
func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	// Two first nodes
	var result bool

	if dhtNode.successor != nil {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(newDHTNode.nodeId))
	}
	//Initiation of ring, 2 node ring
	if dhtNode.successor == nil && dhtNode.predecessor == nil {
		fmt.Println("\n-> changes done to ring: ")
		dhtNode.successor = newDHTNode
		dhtNode.predecessor = newDHTNode
		newDHTNode.successor = dhtNode
		newDHTNode.predecessor = dhtNode
		fmt.Println(dhtNode.nodeId + "-> s " + dhtNode.successor.nodeId)
		fmt.Println(newDHTNode.nodeId + "-> s " + newDHTNode.successor.nodeId + "\n")

		// Connecting last node with first(6 -> 7 -> 0)
	} else if result == true && dhtNode.successor != nil && dhtNode.predecessor != nil {
		dhtNode.successor.predecessor = newDHTNode
		newDHTNode.successor = dhtNode.successor
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		fmt.Println(dhtNode.nodeId + "-> s " + dhtNode.successor.nodeId)
		fmt.Println(newDHTNode.nodeId + "-> s " + newDHTNode.successor.nodeId + "\n")

	} else {
		dhtNode.successor.addToRing(newDHTNode)
	}
}

func (dhtNode *DHTNode) lookup(key string) *DHTNode { /* *DHTNode  */

	if dhtNode.responsible(key) {
		fmt.Println(dhtNode.nodeId)
		return dhtNode
	} else {

		return dhtNode.successor.lookup(key)
	}
}

func (dhtNode *DHTNode) acceleratedLookupUsingFingerTable(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {

	if dhtNode.nodeId == key {
		return true
	} else if dhtNode.predecessor.nodeId == key {
		return false
	}
	result := between([]byte(dhtNode.predecessor.nodeId), []byte(dhtNode.nodeId), []byte(key))
	return result
}

func (dhtNode *DHTNode) printRing() {
	fmt.Println(dhtNode.nodeId)
	printRingHelper(dhtNode, dhtNode.successor)
	fmt.Println("->done\n")

}
func printRingHelper(start *DHTNode, n *DHTNode) {
	if start.nodeId != n.nodeId {
		fmt.Println(n.nodeId)
		printRingHelper(start, n.successor)
	}
}

func (dhtNode *DHTNode) printFingers(m int, bits int) {
	idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("From testCalcFingerTable\nsuccessor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String())
}

func createFingers(dhtNode *DHTNode) [BITS]*DHTNode {
	var nodes [BITS]*DHTNode
	//var distnc [BITS]int

	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		fingerSuccessor := dhtNode.lookup(fingerHex)
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		nodes[i] = fingerSuccessor

	}
	return nodes
}
