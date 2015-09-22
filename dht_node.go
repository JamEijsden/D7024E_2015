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
	sm        chan *Finger
}

func (dhtNode DHTNode) sendMsg(t, address string, wg *sync.WaitGroup) {

	defer wg.Done()
	msg := createMsg(t, dhtNode.nodeId, dhtNode.contact.ip+":"+dhtNode.contact.port, address, dhtNode.contact.ip+":"+dhtNode.contact.port)
	dhtNode.transport.send(msg)
}

/***************** CREATE NODE ***********************************************/

func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	dhtNode.sm = make(chan *Finger)
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

/************************ NODEJOIN *****************************************/
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

/***************************** LOOOOKUUUUPUPP *************************************/

func (dhtNode *DHTNode) lookup(key string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(key) {
		//fmt.Println(dhtNode.nodeId + "is responsible")
		dhtNode.sm <- &Finger{dhtNode.nodeId, src}
	} else {
		// send here
		dhtNode.transport.send(createMsg("lookup", key, src, dhtNode.succ[1], src))
		//return dhtNode.successor.lookup(key)
	}
}

func (dhtNode *DHTNode) lookupForward(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", dhtNode.nodeId, src, msg.Origin, src))
	} else {
		dhtNode.transport.send(createMsg("lookup", msg.Key, src, dhtNode.succ[1], msg.Origin))
	}

}
func (dhtNode *DHTNode) found(f *Finger) {
	dhtNode.sm <- f
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

func (dhtNode *DHTNode) fingerLookup(key string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	fmt.Print("list: \t")
	fmt.Print(dhtNode.fingers.fingerList) //list is empty, wtf?
	if dhtNode.responsible(key) {
		dhtNode.sm <- &Finger{dhtNode.nodeId, src}
	} else {
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length
		for temp >= 0 {
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].hash), []byte(key)) { //check if nodeId and it's last finger is between the key
				if dhtNode.succ[0] == dhtNode.fingers.fingerList[temp].hash && temp == 0 {
					dhtNode.transport.send(createMsg("lookup_finger", key, src, dhtNode.succ[1], src))
				}
				temp = temp - 1
			} else {
				dhtNode.transport.send(createMsg("lookup_finger", key, src, dhtNode.fingers.fingerList[temp].address, src))
			}
		}
	}
}

func (dhtNode *DHTNode) fingerForward(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", msg.Key, src, msg.Origin, src))
	} else {
		fmt.Println("fingerForward: ")
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length
		for temp >= 0 {
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].hash), []byte(msg.Key)) { //check if nodeId and it's last finger is between the key
				if dhtNode.succ[0] == dhtNode.fingers.fingerList[temp].hash && temp == 0 {
					dhtNode.transport.send(createMsg("lookup_finger", msg.Key, src, dhtNode.succ[1], src))
				}
				temp = temp - 1
			} else {
				dhtNode.transport.send(createMsg("lookup_finger", dhtNode.nodeId, src, dhtNode.fingers.fingerList[temp].address, src))
			}
		}
	}
}
