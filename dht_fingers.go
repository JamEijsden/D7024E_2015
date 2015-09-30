package dht

import (
	"encoding/hex"
	"fmt"
	//"sync"
)

const BITS int = 160

type FingerTable struct {
	fingerList [BITS]*Finger
}

type Finger struct {
	hash    string
	address string
}

func fixFingers(dhtNode *DHTNode) {
	if dhtNode.fingers.fingerList[0] == nil {
		dhtNode.fingers = findFingers(dhtNode)
	} else {
		updateFingers(dhtNode)
	}
}

func findFingers(dhtNode *DHTNode) *FingerTable {
	fmt.Print("")
	//fmt.Println(dhtNode)
	var nodes [BITS]*Finger
	var found int
	for i := 0; i < BITS; i++ {
		idBytes, _ := hex.DecodeString(dhtNode.nodeId)
		fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
		//fmt.Println(dhtNode.nodeId + ": FINGERHEX: " + fingerHex)
		if fingerHex == "" {
			fingerHex = "00"
		}
		//if i == 0 {
		//fmt.Println("INITING A LOOkUP")
		go dhtNode.lookup(fingerHex)

		for found != 1 {
			select {
			case s := <-dhtNode.sm:
				//fmt.Println(s)
				if dhtNode.nodeId == s.hash {
					nodes[i] = &Finger{dhtNode.pred[0], dhtNode.pred[1]}
				} else {
					nodes[i] = s
				}
				found = 1
			}
		}
		found = 0
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		//fmt.Printf(fingerSuccessor.nodeId + " ")
		//fmt.Print(nodes[i].nodeId + " ")

	}
	//ft.fingerList = nodes
	return &FingerTable{nodes}
}

func updateFingers(dhtNode *DHTNode) {

	//var nodes [BITS]*Finger
	var found int
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i := 0; i < BITS; i++ {

		if dhtNode.fingers.fingerList[i] != nil {

			idBytes, _ := hex.DecodeString(dhtNode.nodeId)
			fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)

			if i == 0 {
				go dhtNode.lookup(fingerHex)
			} else {
				dhtNode.transport.send(createMsg("lookup", fingerHex, src, dhtNode.fingers.fingerList[i-1].address, src))
			}
			for found != 1 {
				select {
				case s := <-dhtNode.sm:

					found = 1
					dhtNode.fingers.fingerList[i] = s
				}
			}
			found = 0

		}

	}
	/*go func() {
		fmt.Print(dhtNode.nodeId)
		for k := 0; k < BITS; k++ {
			fmt.Print(dhtNode.fingers.fingerList[k])
			fmt.Print(", ")
		}
		fmt.Print("Fingers updated\n")
	}()*/

	//printFingers(dhtNode)

	//ft.fingerList = nodes

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
