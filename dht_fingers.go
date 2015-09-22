package dht

import (
	"encoding/hex"
	"fmt"
)

const BITS int = 3

type FingerTable struct {
	fingerList [BITS]*Finger
}

type Finger struct {
	hash    string
	address string
}

func findFingers(dhtNode *DHTNode) *FingerTable {
	var nodes [BITS]*Finger
	var found int
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		if i == 0 {
			dhtNode.lookup(fingerHex)
		} else {
			dhtNode.transport.send(createMsg("lookup", fingerHex, src, nodes[i-1].address, src))
		}
		for found != 1 {
			select {
			case s := <-dhtNode.sm:
				fmt.Println(s)
				found = 1
				nodes[i] = s
			}
		}
		found = 0
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		//fmt.Printf(fingerSuccessor.nodeId + " ")
		//fmt.Print(nodes[i].nodeId + " ")
	}
	//ft.fingerList = nodes
	fmt.Print("Firstlist yo: ")
	fmt.Println(FingerTable{nodes})
	return &FingerTable{nodes}
}

/*
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
}*/
