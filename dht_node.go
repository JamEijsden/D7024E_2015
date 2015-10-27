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
	listSize    int
	fingerSet   int
	online      int
	nodeId      string
	pred        [2]string
	succ        [2]string
	fingers     *FingerTable
	contact     Contact
	transport   *Transport
	sm          chan *Finger
	responseMsg chan *Msg
	ackMsg      chan *Msg
	taskQueue   chan *Task
	heartbeat   chan string
	succChan    chan *Finger
	dataChannel chan string
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

/*CREATES AND INITIATES NODE */
func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.listSize = 0
	dhtNode.fingerSet = 0
	dhtNode.online = 1
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port
	bindAdr := (dhtNode.contact.ip + ":" + dhtNode.contact.port)
	dhtNode.sm = make(chan *Finger)
	dhtNode.taskQueue = make(chan *Task)
	dhtNode.heartbeat = make(chan string)
	dhtNode.dataChannel = make(chan string)
	dhtNode.responseMsg = make(chan *Msg)
	dhtNode.ackMsg = make(chan *Msg)
	dhtNode.succChan = make(chan *Finger)
	dhtNode.initiated = 0
	if nodeId == nil {
		genNodeId := hashString(bindAdr)
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}
	dhtNode.pred = [2]string{}
	dhtNode.succ = [2]string{dhtNode.nodeId, bindAdr}
	dhtNode.fingers = new(FingerTable)
	dhtNode.fingers.fingerList = [BITS]*Finger{}
	dhtNode.transport = CreateTransport(dhtNode, bindAdr)

	return dhtNode
}

/* Starts an UPD listener */
func (dhtNode *DHTNode) startServer(wg *sync.WaitGroup) {
	dhtNode.online = 1
	go WebServer(dhtNode)
	go initStorage(dhtNode)
	go dhtNode.transport.processMsg()
	go dhtNode.taskHandler()
	go dhtNode.stabilizeTimer()
	go dhtNode.fingerTimer()
	go func() {
		time.Sleep(time.Second * 1)
		go dhtNode.heartbeatHandler()
	}()
	wg.Done()
	dhtNode.transport.listen()

}

/* Prints ring and information of nodes, starting from dhtNode */
func (dhtNode *DHTNode) printRing(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if msg.Type == "" {
		fmt.Println(dhtNode)
		//dhtNode.printFingers()
		fmt.Print("\n")
		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], src))
	} else if src != msg.Origin {
		fmt.Println(dhtNode)
		//dhtNode.printFingers()
		fmt.Print("\n")

		dhtNode.transport.send(createMsg("print", "", src, dhtNode.succ[1], msg.Origin))

	} else {
		fmt.Println("End of ring\n")
	}
}

/* Prints dhtNodes fingers */
func (dhtNode *DHTNode) printFingers() {
	for _, f := range dhtNode.fingers.fingerList {
		fmt.Print(f)
		fmt.Print(" \n")
	}
}

func (dhtNode *DHTNode) gatherAllData(msg *Msg) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if msg.Origin == sender && msg.Src != "" {
		go func() {
			fmt.Println("Im sending data to myself")
			dhtNode.dataChannel <- msg.Key
		}()
	} else {

		msg.Key = msg.Key + getAllData(dhtNode)
		fmt.Println(dhtNode.nodeId + " Im sending my data")
		go dhtNode.transport.send(createMsg("allData", msg.Key, sender, dhtNode.succ[1], msg.Origin))
	}

}

/* Adds dataUpload to Task Queue */
func (dhtNode *DHTNode) queueDataUpload(path string) {
	go dhtNode.QueueTask(createTask("save_data", createMsg("", path, "", "", "")))
}

/* Looks for responsible node for the data and tells it to store it. */
func (dhtNode *DHTNode) storeData(name, data_string string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port

	//data_byte, name := loadData(path)
	hash := hashString(name)

	resp := dhtNode.findSuccesor(hash)
	m := createDataMsg("data_save", name, src, resp.address, src, []byte(data_string))
	go dhtNode.transport.send(m)

}

func (dhtNode *DHTNode) removeData(filename string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port

	//data_byte, name := loadData(path)
	hash := hashString(filename)

	resp := dhtNode.findSuccesor(hash)
	m := createMsg("data_remove", filename, src, resp.address, src)
	go dhtNode.transport.send(m)
}

func (dhtNode *DHTNode) getData(filename string) (string, string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	hash := hashString(filename)
	resp := dhtNode.findSuccesor(hash)
	fmt.Println(resp)
	dhtNode.transport.send(createMsg("request_data", filename, src, resp.address, src))
	for {
		select {
		case d := <-dhtNode.dataChannel:
			//fmt.Println(d)
			return d, resp.hash
		}
	}
}

func (dhtNode *DHTNode) ansDataRequest(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	byte, name := loadData("node_storage/" + dhtNode.nodeId + "/" + msg.Key)
	go dhtNode.transport.send(createDataMsg("data_reply", name, src, msg.Origin, src, byte))
}

/* Checks if Msg.src is dhtNodes predecessor */
func (dhtNode *DHTNode) notify(msg *Msg) {
	if dhtNode.pred[0] == "" || between([]byte(dhtNode.pred[0]), []byte(dhtNode.nodeId), []byte(msg.Key)) {
		dhtNode.pred[0] = msg.Key
		dhtNode.pred[1] = msg.Src
	}
	return
}

/* Checks if successors is correct and fixes if not */
func (dhtNode *DHTNode) stabi() {
	found := false
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	go dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, dhtNode.succ[1], sender))
	waitRespons := time.NewTimer(time.Millisecond * 2000)
	for found != true {
		select {
		case r := <-dhtNode.responseMsg:
			if between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(r.Key)) && r.Key != "" && dhtNode.nodeId != r.Key {
				//fmt.Println(between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(r.Key)) && r.Key != "" && dhtNode.nodeId != r.Key)

				dhtNode.succ[0] = r.Key
				dhtNode.succ[1] = r.Origin
				waitRespons.Stop()
				found = true
				return
			}
			go dhtNode.transport.send(createMsg("notify", dhtNode.nodeId, sender, dhtNode.succ[1], sender))
			return
		case s := <-waitRespons.C:
			s = s
			if dhtNode.online == 1 {
				findAliveSucc(dhtNode)
				fmt.Println("out of findAliveSucc")
				found = true
				return
			}
		}
	}
	found = false
}

func findAliveSucc(dhtNode *DHTNode) {
	waitRespons := time.NewTimer(time.Millisecond * 2)
	waitRespons.Stop()
	found := false
	fmt.Print("findingSucc to: ")
	fmt.Println(dhtNode.nodeId)
	tempAddr := ""
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for i := 1; i < len(dhtNode.fingers.fingerList); i++ {
		if dhtNode.fingers.fingerList[i].address == tempAddr {

		} else {
			waitRespons.Reset(time.Millisecond * 500)
			go dhtNode.transport.send(createMsg("pred", dhtNode.nodeId, sender, dhtNode.fingers.fingerList[i].address, sender))

			for found != true {
				select {
				case r := <-dhtNode.responseMsg:
					fmt.Print("We got a response ")
					fmt.Println(dhtNode.fingers.fingerList[i])
					r = r
					dhtNode.succ[0] = dhtNode.fingers.fingerList[i].hash
					dhtNode.succ[1] = dhtNode.fingers.fingerList[i].address
					waitRespons.Stop()
					found = true
					return

				case s := <-waitRespons.C:
					fmt.Println("Fail response")
					fmt.Println(dhtNode.fingers.fingerList[i].address)
					s = s
					tempAddr = dhtNode.fingers.fingerList[i].address
					found = true
				}
			}
			found = false
		}
	}
}

/* Asks where dhtNode can join and what place in ring */
func (dhtNode *DHTNode) join(adr string) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	dhtNode.pred[0] = ""
	dhtNode.pred[1] = ""
	dhtNode.transport.send(createMsg("findSucc", dhtNode.nodeId, sender, adr, sender))
	for {
		select {
		case s := <-dhtNode.succChan:
			dhtNode.succ[0] = s.hash
			dhtNode.succ[1] = s.address
			go func() {
				time.Sleep(time.Second * 3)
				dhtNode.transport.send(createMsg("update_data", s.hash, sender, dhtNode.succ[1], sender))
			}()
			var nodes [len(dhtNode.fingers.fingerList)]*Finger
			if dhtNode.fingers.fingerList[0] == nil {
				for i := 0; i < len(dhtNode.fingers.fingerList); i++ {
					nodes[i] = &Finger{dhtNode.succ[0], dhtNode.succ[1]}
				}
				dhtNode.fingers = &FingerTable{nodes}
			}
			return
		}
	}
}

/* Finds teh successor of a Hash */
func (dhtNode *DHTNode) startFindSucc(msg *Msg) {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	m := createMsg("foundSucc", dhtNode.succ[0], sender, msg.Origin, dhtNode.succ[1])
	if dhtNode.succ[0] != dhtNode.nodeId && between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key)) {
		dhtNode.transport.send(m)
	} else if dhtNode.nodeId == dhtNode.succ[0] {
		dhtNode.transport.send(m)
	} else {
		dhtNode.transport.send(createMsg("findSucc", msg.Key, msg.Origin, dhtNode.succ[1], sender))
		for {
			select {
			case s := <-dhtNode.succChan:
				m = createMsg("foundSucc", s.hash, sender, msg.Origin, s.address)
				dhtNode.transport.send(m)
				return
			}
		}
	}
	dhtNode.succ[0] = msg.Key
	dhtNode.succ[1] = msg.Origin
}

/* Initiate a successor search */
func (dhtNode *DHTNode) findSuccessorjoin(msg *Msg) {
	result := false
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.succ[0] != dhtNode.nodeId {
		result = between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key))
	} else if dhtNode.succ[0] == dhtNode.nodeId {
		result = true
	}
	if result {
		go dhtNode.transport.send(createMsg("foundSucc", dhtNode.succ[0], sender, msg.Origin, dhtNode.succ[1]))
		dhtNode.succ[0] = msg.Key
		dhtNode.succ[1] = msg.Src
	} else {
		go dhtNode.transport.send(createMsg("findSucc", msg.Key, msg.Src, dhtNode.succ[1], msg.Origin))
	}
}

/* tells a node to change Successor to given successor in msg */
func (dhtNode *DHTNode) reconnNodes(msg *Msg) {
	switch msg.Type {
	case "succ":
		if msg.Src == msg.Origin {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Src
		} else {
			dhtNode.succ[0] = msg.Key
			dhtNode.succ[1] = msg.Origin
		}
		dhtNode.succ[0] = msg.Key
	}
}

/***************************** LOOOOKUUUUPUPP *************************************/
/* finds node responsible for hash*/
func (dhtNode *DHTNode) lookup(key string) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	if dhtNode.responsible(key) {
		dhtNode.sm <- &Finger{key, src}
	} else {
		go dhtNode.transport.send(createMsg("lookup", key, src, dhtNode.succ[1], src))
	}
}

/* forwards lookup request and returns node responsible*/
func (dhtNode *DHTNode) lookupForward(msg *Msg) {
	var waitRespons *time.Timer
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	waitRespons = time.NewTimer(time.Millisecond * 1000)
	if dhtNode.responsible(msg.Key) {
		dhtNode.transport.send(createMsg("lookup_found", dhtNode.nodeId, src, msg.Origin, dhtNode.nodeId))
	} else {
		dhtNode.transport.send(createMsg("lookup", msg.Key, src, dhtNode.succ[1], msg.Origin))
	}
	for {
		select {
		case s := <-dhtNode.ackMsg:
			s = s
			return
		case q := <-waitRespons.C:
			q = q
			return
		}
	}
}

/* saves found node in channel to be taken by other thread waiting */
func (dhtNode *DHTNode) found(f *Finger) {
	dhtNode.sm <- f
}

func (dhtNode *DHTNode) successorFound(succ *Finger) {
	dhtNode.succChan <- succ
}

/* logic for finding responsible node*/
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
				fmt.Print("finger[temp].lookup: ")
				fmt.Println(dhtNode.fingers.fingerList[temp].address)
				dhtNode.transport.send(createMsg("lookup_finger", key, src, dhtNode.fingers.fingerList[temp].address, src))
				temp = -1
			}
		}
	}
	for {
		select {
		case s := <-dhtNode.sm:
			fmt.Println(s.hash)
			return
		}
	}
}

func (dhtNode *DHTNode) startLookup(key string) {
	dhtNode.QueueTask(createTask("fingerLookup", createMsg("", key, "", "", "")))
}

func (dhtNode *DHTNode) findSuccesor(key string) *Finger {
	fmt.Println("Looking up finger: ")
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port

	if between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(key)) {
		return (&Finger{dhtNode.succ[0], dhtNode.succ[1]})
	} else {
		//fmt.Println("sending succ to fingerForward")
		go dhtNode.transport.send(createMsg("lookup_finger", key, src, src, src))
		for {
			select {
			case s := <-dhtNode.succChan:
				fmt.Println("Responisble: " + s.hash + " - " + s.address)
				return s
			}
		}
	}
}

func (dhtNode *DHTNode) fingerForward1(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	//fmt.Println(dhtNode.nodeId + "> " + dhtNode.nodeId + " <-> " + dhtNode.succ[0] + " =? " + msg.Key)

	if between([]byte(dhtNode.nodeId), []byte(dhtNode.succ[0]), []byte(msg.Key)) {
		dhtNode.transport.send(createMsg("foundSucc", dhtNode.succ[0], src, msg.Origin, dhtNode.succ[1]))
		return
	} else {
		length := len(dhtNode.fingers.fingerList) - 1
		temp := length
		for temp >= 0 {
			if between([]byte(dhtNode.nodeId), []byte(msg.Key), []byte(dhtNode.fingers.fingerList[temp].hash)) { //|| dhtNode.fingers.fingerList[temp].hash == msg.Key {

				dhtNode.transport.send(createMsg("lookup_finger", msg.Key, src, dhtNode.fingers.fingerList[temp].address, msg.Origin))
				return
			} else {
				//			fmt.Println("temp - 1 ")
				temp = temp - 1
			}
		}
	}
}

func (dhtNode *DHTNode) stabilize() {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	dhtNode.transport.send(createMsg("stabilize", dhtNode.nodeId, src, dhtNode.succ[1], src))
}

func (dhtNode *DHTNode) stabilizeForward(msg *Msg) {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	go updateFingers(dhtNode)
	if src != msg.Origin {
		dhtNode.transport.send(createMsg("stabilize", dhtNode.nodeId, src, dhtNode.succ[1], msg.Origin))

	}
}

func (dhtNode *DHTNode) QueueTask(t *Task) {
	if dhtNode.online != 1 {
		return
	} else {
		dhtNode.taskQueue <- t
	}
}

func (dhtNode *DHTNode) taskHandler() {
	for {
		if dhtNode.online != 0 {
			select {
			case t := <-dhtNode.taskQueue:
				switch t.Type {
				case "update_data":
					updateData(dhtNode, t.M)

				case "findSucc":
					dhtNode.startFindSucc(t.M)
				case "reconn":
					dhtNode.reconnNodes(t.M)

				case "fixFingers":
					fixFingers(dhtNode)
				case "data":
					dhtNode.storeData(t.M.Key, string(t.M.Data))
					//time.Sleep(time.Millisecond * 500)
				case "notify":
					dhtNode.notify(t.M)
				case "stabilize":
					dhtNode.stabi()
					//dhtNode.stabilizeForward(t.M)
				case "print":
					dhtNode.printRing(t.M)
				case "lookup":
					dhtNode.lookupForward(t.M)
				//case "fingerLookup":
				//	dhtNode.fingerForward(t.M)
				case "fingerLookup":
					dhtNode.findSuccesor(t.M.Key)
				case "leave":
					nMsg := t.M
					if dhtNode.succ[1] == t.M.Src {
						nMsg.Type = "succ"
						dhtNode.reconnNodes(nMsg)
					} else if dhtNode.pred[1] == t.M.Src {
						nMsg.Type = "pred"
						dhtNode.reconnNodes(nMsg)
					}
				case "heartbeat":
					dhtNode.doHeartbeat()
				}
			}
		}
	}
}

func (dhtNode *DHTNode) stabilizeTimer() {
	//src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	for {

		if dhtNode.online == 0 {
			return
		} else {
			time.Sleep(time.Millisecond * 2000)
			go dhtNode.QueueTask(createTask("stabilize", nil))
		}
	}
}

func (dhtNode *DHTNode) fingerTimer() {
	sender := dhtNode.contact.ip + ":" + dhtNode.contact.port
	var nodes [len(dhtNode.fingers.fingerList)]*Finger
	if dhtNode.fingers.fingerList[0] == nil {
		for i := 0; i < len(dhtNode.fingers.fingerList); i++ {
			nodes[i] = &Finger{dhtNode.nodeId, sender}
		}
		dhtNode.fingers = &FingerTable{nodes}
	}

	for {
		if dhtNode.online != 0 {

			time.Sleep(time.Millisecond * 7300)
			go dhtNode.QueueTask(createTask("fixFingers", nil))
			if dhtNode.nodeId == "10" {
				dhtNode.QueueTask(createTask("print", createMsg("", "", "1", "", "")))
			}
		} else {
			return
		}
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
}

/*******************************************************************
******************************* NODE FAIL **************************
********************************************************************/

func (dhtNode *DHTNode) nodeFail() {
	go dhtNode.transport.killConnection()
	return
}

func (dhtNode *DHTNode) heartbeatHandler() {

	for {
		if dhtNode.online != 0 {
			time.Sleep(time.Millisecond * 2000)
			go dhtNode.QueueTask(createTask("heartbeat", nil))

		} else {
			return
		}
	}
}
func (dhtNode *DHTNode) doHeartbeat() {
	src := dhtNode.contact.ip + ":" + dhtNode.contact.port
	respons := false
	state := ""
	var waitRespons *time.Timer
	if dhtNode.pred[0] != "" {
		dhtNode.transport.send(createMsg("heartbeat", dhtNode.nodeId, src, dhtNode.pred[1], src))
		waitRespons = time.NewTimer(time.Millisecond * 1000)
		for !respons {
			select {
			case r := <-dhtNode.heartbeat:
				r = r
				state = "alive"
				respons = true
			case c := <-waitRespons.C:
				c = c
				if dhtNode.online == 1 {
					fmt.Println("This node just died: " + dhtNode.pred[1])
					state = "dead"
					respons = true
				}
			}
		}
		if state == "dead" {
			relocateBackup(dhtNode)
			dhtNode.pred[0] = ""
			dhtNode.pred[1] = ""

		}
		//respons = false
		state = ""
	}
}
