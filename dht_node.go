package dht

import (
	"encoding/hex"
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
	fingers     *FingerTable
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

	dhtNode.fingers = new(FingerTable)
	dhtNode.fingers.fingerList = [BITS]*DHTNode{}
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
		dhtNode.fingers.fingerList = findFingers(dhtNode)
		newDHTNode.fingers.fingerList = findFingers(newDHTNode)
		dhtNode.stabilize(dhtNode.nodeId)
		newDHTNode.stabilize(newDHTNode.nodeId)

		// Connecting last node with first(6 -> 7 -> 0)
	} else if result == true && dhtNode.successor != nil && dhtNode.predecessor != nil {
		dhtNode.successor.predecessor = newDHTNode
		newDHTNode.successor = dhtNode.successor
		dhtNode.successor = newDHTNode
		newDHTNode.predecessor = dhtNode
		fmt.Println(dhtNode.nodeId + "-> s " + dhtNode.successor.nodeId)
		fmt.Println(newDHTNode.nodeId + "-> s " + newDHTNode.successor.nodeId + "\n")
		newDHTNode.fingers.fingerList = findFingers(newDHTNode)
		newDHTNode.stabilize(newDHTNode.nodeId)
	} else {
		dhtNode.successor.addToRing(newDHTNode)
	}

}

func (dhtNode *DHTNode) lookup(key string) *DHTNode { /* *DHTNode  */

	if dhtNode.responsible(key) {
		//fmt.Println(dhtNode.nodeId)
		return dhtNode
	} else {

		return dhtNode.successor.lookup(key)
	}
}

func (dhtNode *DHTNode) acceleratedLookupUsingFingerTable(key string, length int) *DHTNode {
	if dhtNode.responsible(key) {
		return dhtNode
	} else {
		temp := length - 1
		for temp >= 0 {
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].nodeId), []byte(key)) {
				if dhtNode.successor.nodeId == dhtNode.fingers.fingerList[temp].nodeId {
					return dhtNode.fingers.fingerList[temp].acceleratedLookupUsingFingerTable(key, temp)
				}
				temp = temp - 1
			} else {
				return dhtNode.fingers.fingerList[temp].acceleratedLookupUsingFingerTable(key, temp)
			}
		}
	}
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {

	if dhtNode.nodeId == key {
		return true
	} else if dhtNode.predecessor.nodeId == key {
		return false
	}
	isResponsible := between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(key))
	return isResponsible
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

func (dhtNode *DHTNode) stabilize(start string) {
	if dhtNode.successor.nodeId != start {
		fmt.Println(dhtNode.successor.nodeId + " NODELIST: ")
		//	fmt.Println("FUCK MY LIFE IM STOOPID " + dhtNode.successor.nodeId)
		test := updateFingers(dhtNode.successor)
		for i := 0; i < 3; i++ {
			if test[i] != nil {

				fmt.Print(test[i].nodeId + " ")
			}
		}
		fmt.Println("")
		dhtNode.successor.stabilize(start)
	}
	/*
		fmt.Println(dhtNode.nodeId)
		fmt.Print(dhtNode.fingers.fingerList[0].nodeId + " ")
		fmt.Print(dhtNode.fingers.fingerList[1].nodeId + " ")
		fmt.Print(dhtNode.fingers.fingerList[2].nodeId + " \n")
	*/
	/* interate through every node in system
	call func updateFingers() to update current node. */
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
