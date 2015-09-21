package dht

import (
//"encoding/hex"
)

const BITS int = 8

type FingerTable struct {
	fingerList [BITS]*Finger
}

type Finger struct {
	hash    string
	address string
}

/*
func findFingers(dhtNode *DHTNode) *FingerTable {
	var ft = FingerTable{}
	var nodes [BITS]*Finger

	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS) /* returnerar en sträng på vilken nod finger i pengar på.
		fingerSuccessor := dhtNode.lookup(fingerHex)
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		//fmt.Printf(fingerSuccessor.nodeId + " ")
		nodes[i] = fingerSuccessor
		//fmt.Print(nodes[i].nodeId + " ")
	}
	ft.fingerList = nodes
	return ft
}*/

/*
func findFingers(dhtNode *DHTNode) [BITS]*DHTNode {
	var nodes [BITS]*DHTNode /* nodes är en lista en lista med pekare
	//var distnc [BITS]int
	//fmt.Printf(dhtNode.nodeId + " -> Fingers ")
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS) /* returnerar en sträng på vilken nod finger i pengar på.
		fingerSuccessor := dhtNode.lookup(fingerHex)
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		//fmt.Printf(fingerSuccessor.nodeId + " ")
		nodes[i] = fingerSuccessor
		//fmt.Print(nodes[i].nodeId + " ")
	}
	//fmt.Println("Done")
	//fmt.Println("")
	return nodes
}

func updateFingers(dhtNode *DHTNode) [BITS]*DHTNode {

	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		if fingerHex == "" {
			fingerHex = "00"
		}
		if fingerHex != dhtNode.fingers.fingerList[i].nodeId {
			fingerSuccessor := dhtNode.lookup(fingerHex)
			dhtNode.fingers.fingerList[i] = fingerSuccessor

		}
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
	}

	return dhtNode.fingers.fingerList
}
*/
