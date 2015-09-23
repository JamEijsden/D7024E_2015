package dht

import (
	//"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Contact struct {
	ip   string
	port string
}

type Task struct {
	Type string
	M    *Msg
}

type DHTNode struct {
	nodeId    string
	pred      [2]string
	succ      [2]string
	fingers   *FingerTable
	contact   Contact
	transport *Transport
	sm        chan *Finger
	taskQueue chan *Task
	initiated int
}

func (dhtNode DHTNode) sendMsg(t, address string) {

	msg := createMsg(t, dhtNode.nodeId, dhtNode.contact.ip+":"+dhtNode.contact.port, address, dhtNode.contact.ip+":"+dhtNode.contact.port)
	go dhtNode.transport.send(msg)
}

/***************** CREATE NODE ***********************************************/
func createTask(s string, m *Msg) *Task {
	return &Task{s, m}
}
func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	dhtNode.sm = make(chan *Finger)
	dhtNode.taskQueue = make(chan *Task)
	dhtNode.succ = [2]string{}
	dhtNode.pred = [2]string{}
	dhtNode.initiated = 0
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
	go dhtNode.worker()

	return dhtNode
}
func (dhtNode *DHTNode) startServer(wg *sync.WaitGroup) {
	wg.Done()
	dhtNode.transport.listen()

}

func (dhtNode *DHTNode) initRing(msg *Msg) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//func initNode()
	dhtNode.succ[1] = msg.Src
	dhtNode.pred[1] = msg.Src
	dhtNode.succ[0] = msg.Key
	dhtNode.pred[0] = msg.Key
	//fmt.Println(dhtNode)
	fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")

	/*go func() {
		dhtNode.fingers = findFingers(dhtNode)
		//dhtNode.stabilize()
	}()*/
	//fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	/*fmt.Print("initRing msg Type: ")--------
	fmt.Println(msg.Type) // denna är alltid request */
	//Här måste vi ändra msg.Type skulle jag tro. Eller i queueTask skiten.

	if msg.Type == "request" {
		msg := createMsg("init", dhtNode.nodeId, sender, msg.Src, msg.Origin) // RACE CONDITION?
		go dhtNode.transport.send(msg)
	} else if msg.Type == "init" {
		fmt.Println("Ring initiated\n")

		//		dhtNode.QueueTask(createTask("print", nil))
	}

}

func (dhtNode *DHTNode) printRing(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if msg.Type == "" {
		fmt.Println(dhtNode)
		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], src))
	} else if src != msg.Origin {
		fmt.Println(dhtNode)
		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], msg.Origin))

	} else {
		fmt.Println("End of ring")
	}
}

func (dhtNode *DHTNode) joinRing(msg *Msg) {
	var result bool
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.succ[0] != "" {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key))
	}
	/*fmt.Println("We made it to before 2nd if")
	fmt.Println(result)
	fmt.Println(msg.Type)
	fmt.Println(dhtNode.succ[1])
	fmt.Println(dhtNode.pred[1])
	*/ //Connect and reconnect succ/pred when node joins.
	if result == true && dhtNode.succ[1] != "" && dhtNode.pred[1] != "" && msg.Type == "join" {
		dhtNode.transport.send(createMsg("pred", msg.Key, sender, dhtNode.succ[1], msg.Origin))
		dhtNode.transport.send(createMsg("succ", dhtNode.succ[0], sender, msg.Origin, dhtNode.succ[1]))
		//mutex.Lock()
		dhtNode.succ[1] = msg.Origin
		dhtNode.succ[0] = msg.Key
		//mutex.Unlock()
		dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, msg.Origin, sender))
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
		time.Sleep(time.Second * 1)
		//dhtNode.transport.send(createMsg("fingers", dhtNode.nodeId, sender, msg.Origin, sender))
	} else if dhtNode.succ[1] != "" {
		dhtNode.transport.send(createMsg("join", msg.Key, sender, dhtNode.succ[1], msg.Origin))
	} else {
		go dhtNode.QueueTask(createTask("join", msg))
	}
}

/************************ NODEJOIN *****************************************/
func (dhtNode *DHTNode) nodeJoin(msg *Msg) {
	//	var wg sync.WaitGroup
	//Initiate ring if main node doesnt have connections
	if dhtNode.succ[0] == "" && dhtNode.pred[0] == "" && dhtNode.initiated != 1 {
		dhtNode.initiated = 1
		fmt.Println("Beginning init..")
		fmt.Println(msg)
		go dhtNode.QueueTask(createTask("init", msg))

		//KLAR ^
	} else {
		fmt.Println("Beginning join..") //onändlig loop, msg.Type är alltid request i queueTask.
		fmt.Println(msg)
		msg.Type = "join"
		go dhtNode.QueueTask(createTask("join", msg))
	}

	//dhtNode.fingers = findFingers(dhtNode)
}

func (dhtNode *DHTNode) reconnNodes(msg *Msg) {
	switch msg.Type {
	case "pred":
		if msg.Src == msg.Origin {
			dhtNode.pred[0] = msg.Key
			dhtNode.pred[1] = msg.Src
		} else {
			dhtNode.pred[0] = msg.Key
			dhtNode.pred[1] = msg.Origin
		}
		dhtNode.pred[0] = msg.Key

	//	fmt.Println(dhtNode.nodeId + "> Reconnected predecessor, " + msg.Key + "\n")
	case "succ":

		if msg.Src == msg.Origin {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Src
		} else {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Origin
		}
		dhtNode.succ[0] = msg.Key
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

func (dhtNode *DHTNode) stabilize() {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	dhtNode.transport.send(createMsg("stabilize", dhtNode.nodeId, src, dhtNode.succ[1], src))
}

func (dhtNode *DHTNode) stabilizeForward(msg *Msg) {
	//var wg sync.WaitGroup
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//fmt.Println(dhtNode.fingers.fingerList)
	//fmt.Println(dhtNode.succ[1] + " " + msg.Origin)

	go updateFingers(dhtNode)
	if dhtNode.succ[1] != msg.Origin {
		dhtNode.transport.send(createMsg("stabilize", dhtNode.nodeId, src, dhtNode.succ[1], msg.Origin))
	}
}

func (dhtNode *DHTNode) QueueTask(t *Task) {
	/*fmt.Print("Adding ")
	fmt.Print(t.M.Src + " ")
	fmt.Print(t)
	fmt.Println(" to queue")
	*/dhtNode.taskQueue <- t
}

func (dhtNode *DHTNode) worker() {
	for {
		select {
		case t := <-dhtNode.taskQueue:
			switch t.Type {
			case "join":
				//fmt.Println("Exe join")
				if t.M.Type == "request" {
					dhtNode.nodeJoin(t.M)
				} else {
					dhtNode.joinRing(t.M)
				}
			case "init":
				//fmt.Println("Exe init")
				dhtNode.initRing(t.M)
				//time.Sleep(time.Millisecond * 1000)
			case "reconn":
				dhtNode.reconnNodes(t.M)
			case "findFingers":
				fmt.Println("FINGERS!")
			case "print":

				dhtNode.printRing(t.M)
			}
		}
	}
}
