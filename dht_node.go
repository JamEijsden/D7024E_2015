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
	online      int
	nodeId      string
	pred        [2]string
	succ        [2]string
	fingers     *FingerTable
	contact     Contact
	transport   *Transport
	sm          chan *Finger
	responseMsg chan *Msg
	taskQueue   chan *Task
	heartbeat   chan string
	initiated   int
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
	dhtNode.online = 1
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	bindAdr := (dhtNode.contact.ip + ":" + dhtNode.contact.port)
	dhtNode.sm = make(chan *Finger)
	dhtNode.taskQueue = make(chan *Task)
	dhtNode.heartbeat = make(chan string)
	dhtNode.responseMsg = make(chan *Msg)
	dhtNode.initiated = 0
	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.pred = [2]string{}
	dhtNode.succ = [2]string{dhtNode.nodeId, bindAdr}
	dhtNode.fingers = new(FingerTable)
	dhtNode.fingers.fingerList = [BITS]*Finger{}
	dhtNode.transport = CreateTransport(dhtNode, bindAdr)
	go dhtNode.transport.processMsg()
	go dhtNode.taskHandler()
	go dhtNode.updateTimer()

	return dhtNode
}
func (dhtNode *DHTNode) startServer(wg *sync.WaitGroup) {
	wg.Done()
	dhtNode.transport.listen()

}

/*
func (dhtNode *DHTNode) initRing(msg *Msg) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//func initNode()
	dhtNode.succ[1] = msg.Src
	dhtNode.pred[1] = ""
	dhtNode.succ[0] = msg.Key
	dhtNode.pred[0] = ""
	//fmt.Println(dhtNode)
	fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	/*go func() {
		dhtNode.fingers = findFingers(dhtNode)
		//dhtNode.stabilize()
	}()
	//fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	/*fmt.Print("initRing msg Type: ")--------
	fmt.Println(msg.Type) // denna är alltid request
	//Här måste vi ändra msg.Type skulle jag tro. Eller i Task skiten.

	if msg.Type == "request" {
		msg := createMsg("init", dhtNode.nodeId, sender, msg.Src, sender) // RACE CONDITION?
		go dhtNode.transport.send(msg)
		go dhtNode.updateTimer()
		//go dhtNode.heartbeatHandler()
	} else if msg.Type == "init" {
		msg := createMsg("fingers", dhtNode.nodeId, sender, msg.Src, sender)
		go dhtNode.transport.send(msg)
		fmt.Println("Ring initiated\n")
		//go dhtNode.heartbeatHandler()
		time.Sleep(time.Millisecond * 500)
		go dhtNode.QueueTask(createTask("findFingers", nil))
		time.Sleep(time.Millisecond * 500)

		//		dhtNode.QueueTask(createTask("print//", nil))
	}

	//fmt.Println(dhtNode.nodeId + "YO YO FINGER YO")
}*/

func (dhtNode *DHTNode) printRing(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if msg.Type == "" {
		fmt.Println(dhtNode)
		dhtNode.printFingers()
		fmt.Print("\n")
		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], src))
	} else if src != msg.Origin {
		fmt.Println(dhtNode)
		dhtNode.printFingers()
		fmt.Print("\n")
		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], msg.Origin))

	} else {
		fmt.Println("End of ring\n")
	}
}

func (dhtNode *DHTNode) printFingers() {
	for _, f := range dhtNode.fingers.fingerList {
		fmt.Print(f)
		fmt.Print(" \n")

	}
}

func (dhtNode *DHTNode) notify(msg *Msg) {

	if dhtNode.pred[0] == "" || between([]byte(dhtNode.pred[0]), []byte(dhtNode.nodeId), []byte(msg.Key)) {
		dhtNode.pred[0] = msg.Key
		dhtNode.pred[1] = msg.Src
		//fmt.Print(dhtNode.nodeId + "> notify ")
		//fmt.Println(dhtNode)

	}
}

func (dhtNode *DHTNode) stabi() {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, dhtNode.succ[1], sender))
	for {
		select {
		case r := <-dhtNode.responseMsg:
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(r.Key)) && r.Key != "" && dhtNode.nodeId != r.Key {

				dhtNode.succ[0] = r.Key
				dhtNode.succ[1] = r.Origin
			}
			go dhtNode.transport.send(createMsg("notify", dhtNode.nodeId, sender, dhtNode.succ[1], sender))
			fmt.Print(dhtNode.nodeId + ">  ")
			fmt.Print(dhtNode.pred)
			fmt.Print(" - ")
			fmt.Println(dhtNode.succ)
			return
		}
	}

}
func (dhtNode *DHTNode) joinRing(msg *Msg) {
	var result bool
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.succ[0] != dhtNode.nodeId {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key))
	} else if dhtNode.succ[0] == dhtNode.nodeId {
		result = true
	}
	/*fmt.Println("We made it to before 2nd if")
	fmt.Println(result)
	fmt.Println(msg.Type)
	fmt.Println(dhtNode.succ[1])
	fmt.Println(dhtNode.pred[1])
	*/ //Connect and reconnect succ/pred when node joins.
	if result == true && dhtNode.succ[1] != "" && msg.Type == "join" {

		//dhtNode.transport.send(createMsg("pred", msg.Key, sender, dhtNode.succ[1], msg.Origin))

		dhtNode.transport.send(createMsg("succ", dhtNode.succ[0], sender, msg.Origin, dhtNode.succ[1]))
		dhtNode.succ[0] = msg.Key
		dhtNode.succ[1] = msg.Origin
		//dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, msg.Origin, sender))
		//dhtNode.transport.send(createMsg("notify", dhtNode.nodeId, sender, msg.Origin, sender))
		//		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
		//time.Sleep(time.Millisecond * 100)
		//go dhtNode.transport.send(createMsg("fingers", dhtNode.nodeId, sender, msg.Origin, sender))
		go dhtNode.transport.send(createMsg("done", dhtNode.nodeId, sender, msg.Origin, sender))
	} else if dhtNode.succ[1] != "" {
		//fmt.Println(dhtNode)
		dhtNode.transport.send(createMsg("join", msg.Key, sender, dhtNode.succ[1], msg.Origin))
	} else {
		go dhtNode.QueueTask(createTask("join", msg))
	}
}

/************************ NODEJOIN *****************************************/
func (dhtNode *DHTNode) nodeJoin(msg *Msg) {
	//	var wg sync.WaitGroup
	//fmt.Println(msg)
	//Initiate ring if main node doesnt have connections
	/*if dhtNode.succ[0] == "" && dhtNode.pred[0] == "" && dhtNode.initiated != 1 {
		dhtNode.initiated = 1
		//	fmt.Println("Beginning init..")
		//	fmt.Println(msg)
		go dhtNode.QueueTask(createTask("init", msg))

		//KLAR
	}*/if msg.Type == "done" {
		//go dhtNode.heartbeatHandler()
		fmt.Println(dhtNode.succ[0] + "<- succ :" + dhtNode.nodeId + ": pred -> " + dhtNode.pred[0] + "\n")
	} else {
		//fmt.Println("Beginning join..") //onändlig loop, msg.Type är alltid request i queueTask.
		msg.Type = "join"
		go dhtNode.QueueTask(createTask("join", msg))
		//}
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
		go dhtNode.transport.send(createMsg("lookup", key, src, dhtNode.succ[1], src))
		//return dhtNode.successor.lookup(key)
	}
}

func (dhtNode *DHTNode) lookupForward(msg *Msg) {
	//fmt.Println(dhtNode.nodeId)
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", dhtNode.nodeId, src, msg.Origin, src))
	} else {
		//fmt.Println(msg)
		dhtNode.transport.send(createMsg("lookup", msg.Key, src, dhtNode.succ[1], msg.Origin))
	}

}
func (dhtNode *DHTNode) found(f *Finger) {
	dhtNode.sm <- f

}

func (dhtNode *DHTNode) responsible(key string) bool {

	if dhtNode.nodeId == key {
		//fmt.Println(dhtNode.nodeId + ": I R RESPONSIBLE FOR " + key)
		return true
	} else if dhtNode.pred[0] == key {
		//fmt.Println(dhtNode.nodeId + ": PRED R RESPONSIBLE FOR " + key)
		return false
	}
	isResponsible := between([]byte(dhtNode.pred[0]), []byte(dhtNode.nodeId), []byte(key))
	return isResponsible
}

func (dhtNode *DHTNode) fingerLookup(key string) {
	fmt.Println("Looking up finger: ")
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	found := 0
	if dhtNode.responsible(key) {
		go dhtNode.found(&Finger{dhtNode.nodeId, src})
	} else {
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length
		for temp >= 0 {
			fmt.Print(dhtNode.nodeId)
			fmt.Print(" - ")
			fmt.Print(dhtNode.fingers.fingerList[temp].hash)
			fmt.Println("")
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].hash), []byte(key)) {
				temp = temp - 1
			} else {
				dhtNode.transport.send(createMsg("lookup_finger", key, src, dhtNode.fingers.fingerList[temp].address, src))
				temp = -1
			}
		}
	}
	for found != 1 {
		select {
		case s := <-dhtNode.sm:
			fmt.Println(s.hash)
			found = 1
		}
	}
}

func (dhtNode *DHTNode) fingerForward(msg *Msg) {
	//fmt.Println(dhtNode.nodeId)
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", msg.Key, src, msg.Origin, src))
	} else {
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length
		for temp >= 0 {
			fmt.Print(dhtNode.nodeId)
			fmt.Print(" - ")
			fmt.Print(dhtNode.fingers.fingerList[temp].hash)
			fmt.Println("")
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.fingers.fingerList[temp].hash), []byte(msg.Key)) { //check if nodeId and it's last finger is between the key
				temp = temp - 1
			} else {
				dhtNode.transport.send(createMsg("lookup_finger", msg.Key, src, dhtNode.fingers.fingerList[temp].address, msg.Origin))
				temp = -1
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
	if src != msg.Origin {
		dhtNode.transport.send(createMsg("stabilize", dhtNode.nodeId, src, dhtNode.succ[1], msg.Origin))
	}
}

func (dhtNode *DHTNode) QueueTask(t *Task) {
	/*fmt.Print("Adding ")
	fmt.Print(dhtNode.nodeId + " ")
	fmt.Print(t)
	fmt.Println(" to queue")
	*/dhtNode.taskQueue <- t
}

func (dhtNode *DHTNode) taskHandler() {
	for {
		if dhtNode.online != 0 {
			select {
			case t := <-dhtNode.taskQueue:
				switch t.Type {
				case "join":
					//fmt.Println("Exe join")
					if t.M.Type == "request" || t.M.Type == "done" {
						dhtNode.nodeJoin(t.M)
					} else {
						dhtNode.joinRing(t.M)
					}
				case "init":
					//fmt.Println("Exe init")
					//dhtNode.initRing(t.M)
					//time.Sleep(time.Millisecond * 1000)

				case "reconn":
					dhtNode.reconnNodes(t.M)
				case "findFingers":
					dhtNode.fingers = findFingers(dhtNode)
					//time.Sleep(time.Millisecond * 500)
				case "notify":
					//dhtNode.stabilize()
					dhtNode.notify(t.M)
				case "stabilize":
					dhtNode.stabi()
					//dhtNode.stabilizeForward(t.M)
				case "print":
					dhtNode.printRing(t.M)
				case "lookup":
					dhtNode.lookupForward(t.M)
				case "fingerLookup":
					dhtNode.fingerForward(t.M)
				case "leave":
					nMsg := t.M
					if dhtNode.succ[1] == t.M.Src {
						nMsg.Type = "succ"
						dhtNode.reconnNodes(nMsg)
					} else if dhtNode.pred[1] == t.M.Src {
						nMsg.Type = "pred"
						dhtNode.reconnNodes(nMsg)
					}
				}
			}
		}
	}
}

func (dhtNode *DHTNode) updateTimer() {
	//src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for {
		time.Sleep(time.Millisecond * 5000)
		//dhtNode.transport.send(createMsg("notify", dhtNode.nodeId, src, dhtNode.succ[1], src))
		go func() {
			dhtNode.QueueTask(createTask("stabilize", nil))
		}()
		//dhtNode.QueueTask(createTask("print", createMsg("", "", "1", "", "")))
		//go dhtNode.QueueTask(createTask("startStabilize", nil))
	}

}

/******************************************************************
***************************** NODE LEAVE **************************
*******************************************************************/

func (dhtNode *DHTNode) leaveRing() {
	dhtNode.online = 0
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	go dhtNode.transport.send(createMsg("leave", dhtNode.pred[0], src, dhtNode.succ[1], dhtNode.pred[1]))
	go dhtNode.transport.send(createMsg("leave", dhtNode.succ[0], src, dhtNode.pred[1], dhtNode.succ[1]))
	fmt.Println(dhtNode.nodeId + " left the ring")
	/*	go func() {
			for true {
				fmt.Println("Tomte")
			}
		}()
	*/
}

/*******************************************************************
******************************* NODE FAIL **************************
********************************************************************/

func (dhtNode *DHTNode) nodeFail() {
	dhtNode.online = 0
}

func (dhtNode *DHTNode) heartbeatHandler() {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	respons := false
	state := ""
	var waitRespons *time.Timer
	for {
		if dhtNode.online != 0 {
			time.Sleep(time.Millisecond * 2000)
			go dhtNode.transport.send(createMsg("heartbeat", dhtNode.nodeId, src, dhtNode.pred[1], src))

			waitRespons = time.NewTimer(time.Millisecond * 200)
			for !respons {
				select {
				case r := <-dhtNode.heartbeat:
					r = r
					fmt.Println(dhtNode.nodeId + "> " + r + " is alive!")
					respons = true
					state = "alive"
				case c := <-waitRespons.C:
					//fmt.Println(c)
					c = c
					fmt.Println(dhtNode.nodeId + "> Didnt recieve respons")
					respons = true
					state = "dead"
				}
			}
			if state == "dead" {
				fmt.Println("go dhtNode.QueueTask()")
				//dhtNode.transport.send(createMsg("fail", "", s, d, o))

			} else {
				//fmt.Println(state)
			}
			respons = false
			state = ""
		}
	}
}
