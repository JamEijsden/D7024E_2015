package dht

import (
	//"encoding/hex"
	"fmt"
	"sync"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId    string
	pred      [2]string
	succ      [2]string
	fingers   *FingerTable
	contact   Contact
	transport *Transport
}

func (dhtNode DHTNode) sendMsg(t, address string, wg *sync.WaitGroup) {

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

	dhtNode.fingers = new(FingerTable)
	dhtNode.fingers.fingerList = [BITS]*Finger{}
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
	//Initiate ring if main node doesnt have connections
	if dhtNode.succ[0] == "" && dhtNode.pred[0] == "" {
		dhtNode.succ[1] = msg.Src
		dhtNode.pred[1] = msg.Src
		dhtNode.succ[0] = msg.Key
		dhtNode.pred[0] = msg.Key
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
		if msg.Type == "init" {
			msg := createMsg("ack", dhtNode.nodeId, sender, msg.Src, msg.Origin) // RACE CONDITION?
			dhtNode.transport.send(msg)
		} else if msg.Type == "ack" {
			fmt.Println("Ring initiated")

		}
		//Connect and reconnect succ/pred when node joins.
	} else if result == true && dhtNode.succ[1] != "" && dhtNode.pred[1] != "" && msg.Type == "join" {
		dhtNode.transport.send(createMsg("pred", msg.Key, sender, dhtNode.succ[1], msg.Origin))

		dhtNode.transport.send(createMsg("succ", dhtNode.succ[0], sender, msg.Origin, dhtNode.succ[1]))
		dhtNode.succ[1] = msg.Origin
		dhtNode.succ[0] = msg.Key
		dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, msg.Origin, sender))
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	} else {
		dhtNode.transport.send(createMsg("join", msg.Key, sender, dhtNode.succ[1], msg.Origin))
	}
}

func (dhtNode *DHTNode) reconnNodes(msg *Msg) {
	mutex := &sync.Mutex{}
	switch msg.Type {
	case "pred":
		mutex.Lock()
		if msg.Src == msg.Origin {
			dhtNode.pred[0] = msg.Key
			dhtNode.pred[1] = msg.Src
		} else {
			dhtNode.pred[0] = msg.Key
			dhtNode.pred[1] = msg.Origin
		}
		dhtNode.pred[0] = msg.Key
		mutex.Unlock()
	//	fmt.Println(dhtNode.nodeId + "> Reconnected predecessor, " + msg.Key + "\n")
	case "succ":
		mutex.Lock()
		if msg.Src == msg.Origin {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Src
		} else {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Origin
		}
		dhtNode.succ[0] = msg.Key
		mutex.Unlock()
		//	fmt.Println(dhtNode.nodeId + "> Reconnected successor, " + msg.Key + "\n")
	}
	//fmt.Println(dhtNode)
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

func (dhtNode *DHTNode) lookup(key string) *Finger {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(key) {
		//fmt.Println(dhtNode.nodeId)
		return &Finger{dhtNode.nodeId, src}
	} else {
		// send here
		dhtNode.transport.send(createMsg("lookup", key, src, dhtNode.succ[1], src))
		//return dhtNode.successor.lookup(key)
		return dhtNode.fingers.fingerList[0]
	}
}

func (dhtNode *DHTNode) lookupForward(msg *Msg) {
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", k, s, d, o))
	}

}

func (dhtNode *DHTNode) responsible(key string) bool {

	if dhtNode.nodeId == key {
		return true
	} else if dhtNode.pred[0] == key {
		return false
	}
	isResponsible := between([]byte(dhtNode.pred[0]), []byte(dhtNode.nodeId), []byte(key))
	return isResponsible
}

/***************************************************** WITHOUT NETWORK, ONLY LOCAL *********************************************/
