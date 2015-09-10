package dht

import (
	"encoding/hex"
	"fmt"
	"sync"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	pred        [2]string
	succ        [2]string
	fingers     *FingerTable
	contact     Contact
	transport   *Transport
}

func (dhtNode DHTNode) helloWorld(t, address string, wg *sync.WaitGroup) {

	defer wg.Done()
	msg := createMsg(t, dhtNode.nodeId, dhtNode.contact.ip+":"+dhtNode.contact.port, address, dhtNode.contact.ip+":"+dhtNode.contact.port)
	dhtNode.transport.send(msg)
}

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	dhtNode.succ = [2]string{}
	dhtNode.pred = [2]string{}
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
	bindAdr := (dhtNode.contact.ip + ":" + dhtNode.contact.port)
	dhtNode.transport = CreateTransport(dhtNode, bindAdr)
	go dhtNode.transport.processMsg()

	return dhtNode
}
func (dhtNode *DHTNode) startServer(wg *sync.WaitGroup) {
	wg.Done()
	dhtNode.transport.listen()

}

func (dhtNode *DHTNode) nodeJoin(msg *Msg) {
	//var wg sync.WaitGroup
	var result bool
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.succ[0] != "" {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key))
	}
	if dhtNode.succ[0] == "" && dhtNode.pred[0] == "" {
		dhtNode.succ[1] = msg.Src
		dhtNode.pred[1] = msg.Src
		dhtNode.succ[0] = msg.Key
		dhtNode.pred[0] = msg.Key
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
		if msg.Type != "init" {
			msg := createInitMsg(dhtNode.nodeId, sender, msg.Src) // RACE CONDITION?
			dhtNode.transport.send(msg)
		} else if msg.Type == "init" {
			fmt.Println("Ring initiated")

		}
	} else if result == true && dhtNode.succ[1] != "" && dhtNode.pred[1] != "" {

		go dhtNode.transport.send(createMsg("pred", msg.Key, sender, dhtNode.succ[1], msg.Origin))
		go dhtNode.transport.send(createMsg("succ", dhtNode.succ[0], sender, msg.Src, sender))
		dhtNode.succ[1] = msg.Src
		dhtNode.succ[0] = msg.Key
		go dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, msg.Src, sender))
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	} else {
		go dhtNode.transport.send(createMsg("join", msg.Key, sender, dhtNode.succ[1], msg.Origin))
	}
}

func (dhtNode *DHTNode) reconnNodes(msg *Msg) {
	mutex := &sync.Mutex{}
	switch msg.Type {
	case "pred":
		mutex.Lock()
		if msg.Src == msg.Src {
			dhtNode.pred[1] = msg.Src
		} else {
			dhtNode.pred[1] = msg.Origin
		}
		dhtNode.pred[0] = msg.Key
		mutex.Unlock()
		fmt.Println(dhtNode.nodeId + "> Reconnected predecessor, " + msg.Key + "\n")
	case "succ":
		mutex.Lock()
		if msg.Src == msg.Src {
			dhtNode.succ[1] = msg.Src
		} else {
			dhtNode.succ[1] = msg.Origin
		}
		dhtNode.succ[0] = msg.Key
		mutex.Unlock()
		fmt.Println(dhtNode.nodeId + "> Reconnected successor, " + msg.Key + "\n")
	}
}

func (dhtNode *DHTNode) netPrintRing(msg *Msg) {
	var orig string
	fmt.Println(dhtNode.nodeId + ">" + dhtNode.contact.port)
	if msg == nil {
		orig = dhtNode.contact.ip + ":" + dhtNode.contact.port
		go dhtNode.transport.send(createMsg("circle", "", dhtNode.contact.ip+":"+dhtNode.contact.port, dhtNode.succ[1], orig))
	}
	if dhtNode.succ[1] != msg.Origin && msg != nil {
		go dhtNode.transport.send(createMsg("circle", "", dhtNode.contact.ip+":"+dhtNode.contact.port, dhtNode.succ[1], msg.Origin))
	}
}

/***************************************************** WITHOUT NETWORK, ONLY LOCAL *********************************************/

func testBetween(id1, id2, key string) {
	fmt.Println(id1 + " " + id2 + " " + key + " ")
	fmt.Println(between([]byte(id1), []byte(id2), []byte(key)))
}

// Connects and rearranges nodes
func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	// Two first nodes
	var result bool

	if dhtNode.successor != nil {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.successor.nodeId), []byte(newDHTNode.nodeId))
	}
	//Initiation of ring, 2 node ring
	if dhtNode.successor == nil && dhtNode.predecessor == nil {
		//fmt.Println("\n-> changes done to ring: ")
		dhtNode.successor = newDHTNode
		dhtNode.predecessor = newDHTNode
		newDHTNode.successor = dhtNode
		newDHTNode.predecessor = dhtNode
		//fmt.Println(dhtNode.nodeId + "-> s " + dhtNode.successor.nodeId)
		//fmt.Println(newDHTNode.nodeId + "-> s " + newDHTNode.successor.nodeId + "\n")
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
		//fmt.Println(dhtNode.nodeId + "-> s " + dhtNode.successor.nodeId)
		//fmt.Println(newDHTNode.nodeId + "-> s " + newDHTNode.successor.nodeId + "\n")
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

func (dhtNode *DHTNode) acceleratedLookupUsingFingerTable(key string) *DHTNode {
	if dhtNode.responsible(key) {
		return dhtNode
	} else {
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length //length of list
		for temp >= 0 {
			//fmt.Println(dhtNode.nodeId + " - finger: " + dhtNode.fingers.fingerList[temp].nodeId)
			//fmt.Println(temp)
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].nodeId), []byte(key)) { //check if nodeId and it's last finger is between the key
				if dhtNode.successor.nodeId == dhtNode.fingers.fingerList[temp].nodeId && temp == 0 {
					return dhtNode.successor
				}
				temp = temp - 1
				//fmt.Println(temp)
			} else {
				//fmt.Println("change node: " + dhtNode.fingers.fingerList[temp].nodeId)
				return dhtNode.fingers.fingerList[temp].acceleratedLookupUsingFingerTable(key)
			}
		}
	}
	return dhtNode
}

func (dhtNode *DHTNode) responsible(key string) bool {

	if dhtNode.nodeId == key {
		return true
	} else if dhtNode.predecessor.nodeId == key {
		return false
	}
	isResponsible := between([]byte(dhtNode.predecessor.nodeId), []byte(dhtNode.nodeId), []byte(key))
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
		//	fmt.Println(dhtNode.successor.nodeId + " NODELIST: ")
		//	fmt.Println("FUCK MY LIFE IM STOOPID " + dhtNode.successor.nodeId)
		test := updateFingers(dhtNode.successor)
		for i := 0; i < 3; i++ {
			if test[i] != nil {

				//			fmt.Print(test[i].nodeId + " ")
			}
		}
		//	fmt.Println("")
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
	//fmt.Println("From testCalcFingerTable\nsuccessor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String())
}
