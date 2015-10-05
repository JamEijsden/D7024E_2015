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
						//go transport.node.QueueTask(createTask("lookup", m))
						go transport.node.lookupForward(m)
					//case "lookup_finger_again":
					//	go transport.node.findSuccesor(m.Key)
					case "lookup_finger":
						go transport.node.fingerForward1(m)
					case "lookup_found":
						//fmt.Print("found this: ")
						//fmt.Println(&Finger{m.Key, m.Src})
						//fmt.Println(transport.node.nodeId + "YO YO FOUND YO")
						go transport.node.found(&Finger{m.Key, m.Src})
						//fmt.Println(&Finger{m.Key, m.Src})
					case "notify":
						go transport.node.QueueTask(createTask("notify", m))
					case "stabilize":
						go transport.node.QueueTask(createTask("stabilize", m))
						//transport.node.stabilizeForward(m)
					case "print":
						go transport.node.QueueTask(createTask("print", m))
					case "leave":
						go transport.node.QueueTask(createTask("leave", m))
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

					case "data":
						go replicate(transport.node, m)
					}
				}
			}
		}
	}()
}

func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress) //adds adress to variable udpAddr and err msg in var err is there is one.

	transport.conn, err = net.ListenUDP("udp", udpAddr) //we listen to the IP-adress in udpAddr.
	//fmt.Println("Server running on " + transport.bindAddress + " with ID " + transport.node.nodeId + "\nWaiting for messages..")

	if err != nil {
		fmt.Println(err.Error())
	}
	defer transport.conn.Close()
	dec := json.NewDecoder(transport.conn) //decodes conn and adds to variable dec
	for {
		if transport.node.online == 1 {
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
}

func (transport *Transport) killConnection() {
	transport.conn.Close()
}

func (transport *Transport) send(msg *Msg) {
	//fmt.Println(msg.Src + "> Type: " + msg.Key + " " + msg.Type + " to " + msg.Dst)
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	if err != nil {
		fmt.Println(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
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
