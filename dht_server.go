package dht

import (
	"encoding/json"
	"fmt"
	"net"
)

type Transport struct {
	node        *DHTNode
	bindAddress string
	queue       chan *Msg
	conn        *net.UDPConn
}

func CreateTransport(node *DHTNode, bindAddress string) *Transport {
	transport := &Transport{}
	transport.queue = make(chan *Msg)
	transport.node = node
	transport.bindAddress = bindAddress
	return transport
}

func (transport *Transport) processMsg() {
	//msg := <-transport.queue
	//joinChan := make(chan *Msg)

	go func() {
		for {
			if transport.node.online != 0 {
				select {
				case m := <-transport.queue:
					switch m.Type {
					case "update_data":
						go transport.node.QueueTask(createTask("update_data", m))
					case "init":
						//fmt.Println(transport.node.nodeId + " INITING")
						go transport.node.QueueTask(createTask("init", m))
					case "done":
						//fmt.Println(m.Key)
						//go transport.node.QueueTask(createTask("join", m))
					case "pred":
						src := transport.node.contact.ip + ":" + transport.node.contact.port
						go transport.send(createMsg("response", transport.node.pred[0], src, m.Origin, transport.node.pred[1]))
					case "foundSucc":
						go transport.node.successorFound(&Finger{m.Key, m.Origin})
					case "succ":
						go transport.node.successorFound(&Finger{m.Key, m.Origin})
						//transport.node.QueueTask(createTask("reconn", m))
					case "lookup":
						go transport.send(createMsg("ack", "", m.Src, m.Src, m.Src))
						go transport.node.lookupForward(m)
						//go transport.node.QueueTask(createTask("lookup", m))
					//case "lookup_finger_again":
					//	go transport.node.findSuccesor(m.Key)
					case "lookup_finger":
						go transport.send(createMsg("response", "", m.Src, m.Src, m.Src))
						go transport.node.fingerForward1(m)
					case "lookup_found":
						go transport.send(createMsg("ack", "", m.Src, m.Src, m.Src))
						if m.Key == m.Origin {
							go transport.node.found(&Finger{m.Key, m.Src}) //fmt.Println(&Finger{m.Key, m.Src})
						} else {
							go transport.node.found(&Finger{m.Key, m.Src})
						}
					case "notify":
						//fmt.Println("Stabi notify -> notify server")
						go transport.node.QueueTask(createTask("notify", m))
					case "stabilize":
						go transport.node.QueueTask(createTask("stabilize", m))
						//transport.node.stabilizeForward(m)
					case "print":
						go transport.node.QueueTask(createTask("print", m))
					case "leave":
						go transport.node.QueueTask(createTask("leave", m))

					//TESTAR EN SAK
					case "ack":
						go func() {
							transport.node.ackMsg <- m
						}()
					case "response":
						//go transport.node.QueueTask(createTask("stabilize", m))
						go func() {
							transport.node.responseMsg <- m
						}()
					case "heartbeat":
						go transport.node.heartbeatRespons(m)

						//go transport.node.QueueTask(createTask("foundSucc", m))
					case "findSucc":
						//go transport.node.findSuccessorjoin(msg)
						if m.Src == m.Origin {
							go transport.node.QueueTask(createTask("findSucc", m))
						} else {
							go transport.node.findSuccessorjoin(m)

						}
					case "heartbeat_respons":
						go func() {
							transport.node.heartbeat <- m.Key
						}()

<<<<<<< HEAD
					case "data":
						if m.Src != transport.node.pred[1] || m.Origin == transport.node.pred[1] {
							go savedata(transport.node, m)
						} else {
							go replicate(transport.node, m)
						}

=======
					case "data_save":
						go replicate(transport.node, m)
					case "data_remove":
						go deleteData(transport.node, m)
>>>>>>> 4f5493a616e2e83ec65d8cca33b3b40a0b04b1c8
					case "allData":
						go transport.node.gatherAllData(m)
					case "request_data":
						go transport.node.ansDataRequest(m)
					case "data_reply":
						go func() {
							transport.node.dataChannel <- string(m.Data)
						}()
					}
				}
			} else {
				return
			}
		}
	}()
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress) //adds adress to variable udpAddr and err msg in var err is there is one.

	transport.conn, err = net.ListenUDP("udp", udpAddr)
	transport.conn.SetReadBuffer(10000)
	transport.conn.SetWriteBuffer(10000) //we listen to the IP-adress in udpAddr.
	//fmt.Println("Server running on " + transport.bindAddress + " with ID " + transport.node.nodeId + "\nWaiting for messages..")

	if err != nil {
		fmt.Println(err.Error())
	}
	defer transport.conn.Close()
	dec := json.NewDecoder(transport.conn) //decodes conn and adds to variable dec
	for {
		if transport.node.online == 0 {
			return
		}
		msg := Msg{}
		err = dec.Decode(&msg) //decodes the message where the message adress is and adds an error msg if there's none.
		//fmt.Println(transport.bindAddress + "> Received Message from " + msg.Src)
		//fmt.Println(msg)
		go func() {
			transport.queue <- &msg
		}()
		//return
		//	we	got	a	message, do something

	}
}

func (transport *Transport) killConnection() {
	transport.node.online = 0
	transport.conn.Close()
	return
}

func (transport *Transport) send(msg *Msg) {
	//fmt.Println(msg.Src + "> Type: " + msg.Key + " " + msg.Type + " to " + msg.Dst)
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	conn.SetWriteBuffer(10000)
	conn.SetReadBuffer(10000)
	if err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Println(msg)
	bytes, err := json.Marshal(msg)

	defer conn.Close()
	_, err = conn.Write(bytes)
	if err != nil {
		fmt.Println(err.Error())
	}

}
