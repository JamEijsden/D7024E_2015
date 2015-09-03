package dht

import (
	"encoding/hex"
	"fmt"
)

const BITS int = 3

type FingerTable struct {
	fingerList [BITS]*DHTNode
}

func findFingers(dhtNode *DHTNode) [BITS]*DHTNode {
	var nodes [BITS]*DHTNode
	//var distnc [BITS]int
	fmt.Printf("[")
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		fingerSuccessor := dhtNode.lookup(fingerHex)
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		fmt.Printf(fingerSuccessor.nodeId + " ")
		nodes[i] = fingerSuccessor

	}
	fmt.Printf("]\n")
	return nodes
}

func (dhtNode *DHTNode) updateFingers() {
	fmt.Println("stub")
}
