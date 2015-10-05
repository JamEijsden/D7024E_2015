package dht

import (
	"time"
)

type Msg struct {
	Timestamp int64 `json:"time"`

	Type string `json:"type"`

	Key string `json:"key"`

	Src string `json:"src"`

	Dst string `json:"dst"`

	Origin string `json:"origin"`

	Data []byte `json:"data"`
	//Node *DHTNode `json:"node"`
}

func createMsg(t, k, s, d, o string) *Msg {
	msg := new(Msg)
	msg.Type = t
	msg.Key = k
	msg.Src = s
	msg.Dst = d
	msg.Origin = o
	//msg.Node = n
	msg.Timestamp = time.Now().UnixNano()
	return msg
}

func createDataMsg(t, k, s, d, o string, data []byte) *Msg {
	msg := createMsg(t, k, s, d, o)
	msg.Data = data
	msg.Timestamp = time.Now().UnixNano()
	return msg
}
