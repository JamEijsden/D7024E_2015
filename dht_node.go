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

func (dhtNode *DHTNode) lookup(key string) { /* *DHTNode */
	node := lookuphelp(dhtNode, key)
	if node == dhtNode {
		fmt.Println("LOOPEEDOOPEDOOPEEDOOOOOOOP")
	}
	fmt.Println("nodeID = " + node.nodeId)
}

func (dhtNode *DHTNode) lookuphelp(start *DHTNode, key string) *DHTNode {
	if dhtNode.nodeId != start.nodeId {
		result := between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(key)) /* tar nodeX, nodeX+1 och en hash jag kollar upp */
		if result {
			return dhtNode
		} else {
			dhtNode.successor.lookuphelp(start, key) //infinite loop?
		}
	}
	return dhtNode /* We've done a loop. */
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
	fmt.Println(dhtNode.nodeId)
	printRingHelper(dhtNode, dhtNode.successor)

}
func printRingHelper(start *DHTNode, n *DHTNode) {
	if start.nodeId != n.nodeId {
		fmt.Println(n.nodeId)
		printRingHelper(start, n.successor)
	}
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
