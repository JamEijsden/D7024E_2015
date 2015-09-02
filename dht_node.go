package dht

import (
	"fmt"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = nil
	dhtNode.predecessor = nil

	return dhtNode
}

func testBetween(id1, id2, key string) {
	fmt.Println(id1 + " " + id2 + " " + key + " ")
	fmt.Println(between([]byte(id1), []byte(id2), []byte(key)))
}

func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	// Two first nodes
	var result bool

	if dhtNode.successor != nil {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(newDHTNode.nodeId))
	}

	if dhtNode.successor == nil && dhtNode.predecessor == nil {
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

func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}

func (dhtNode *DHTNode) printRing() {
	fmt.Println("Done")
}

func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
	/* idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String()) */
}
