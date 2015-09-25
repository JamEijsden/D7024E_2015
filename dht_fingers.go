package dht

import (
	"encoding/hex"
	"fmt"
	"sync"
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
	fmt.Print("")
	//fmt.Println(dhtNode)
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
		fmt.Println("WABALUBADUBDUB")
	}
	//ft.fingerList = nodes
	return &FingerTable{nodes}
}

func updateFingers(dhtNode *DHTNode) *FingerTable {
	var m = &sync.Mutex{}
	var nodes [BITS]*Finger
	var found int
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i := 0; i < BITS; i++ {
		m.Lock()
		// defer m.Unlock()
		if dhtNode.fingers.fingerList[i] != nil {

			idBytes, _ := hex.DecodeString(dhtNode.nodeId)
			fingerHex, _ := calcFinger(idBytes, (i + 1), BITS)
			//IF finger havent changed, dont bother lookup
			//fmt.Println(dhtNode.fingers.fingerList[i])
			if fingerHex != dhtNode.fingers.fingerList[i].hash {
				if i == 0 {
					dhtNode.lookup(fingerHex)
				} else {
					dhtNode.transport.send(createMsg("lookup", fingerHex, src, dhtNode.fingers.fingerList[i-1].address, src))
				}
				for found != 1 {
					select {
					case s := <-dhtNode.sm:
						//fmt.Println(s)
						//		fmt.Print(dhtNode.nodeId + "> " + fingerHex + " =? ")
						//		fmt.Print(s)
						//		fmt.Println("")
						found = 1
						dhtNode.fingers.fingerList[i] = s
					}
				}
				found = 0
			} else {
				//fmt.Println(fingerHex + " =keep? " + dhtNode.fingers.fingerList[i].hash)

			}
		}
		m.Unlock()
		//fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
		//dist := distance(idBytes, fingerSuccessorBytes, BITS)
		//fmt.Printf(fingerSuccessor.nodeId + " ")
		//fmt.Print(nodes[i].nodeId + " ")
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
