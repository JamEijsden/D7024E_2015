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
	var nodes [BITS]*DHTNode /* nodes är en lista en lista med pekare */
	//var distnc [BITS]int
	fmt.Printf("[")
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS) /* returnerar en sträng på vilken nod finger i pengar på. */
		fingerSuccessor := dhtNode.lookup(fingerHex)
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		fmt.Printf(fingerSuccessor.nodeId + " ")
		nodes[i] = fingerSuccessor

	}
	fmt.Printf("]\n")
	return nodes
}

func updateFingers(dhtNode *DHTNode) [BITS]*DHTNode {
	var nodes [BITS]*dhtNode
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		if fingerHex != dhtNode.fingers.fingerList[i] {
			fingerSuccessor := dhtNode.lookup(fingerHex)
			nodes[i] = fingerSuccessor
		}

		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)

	}
	return nodes
}
