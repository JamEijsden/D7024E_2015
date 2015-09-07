package dht

import (
	"time"
)

type Msg struct {
	Timestamp int64

	Type string

	Key string

	Src string

	Dst string

	Origin string
}

func createMsg(t, k, s, d, o string) *Msg {
	msg := new(Msg)
	msg.Type = t
	msg.Key = k
	msg.Src = s
	msg.Dst = d
	msg.Origin = o
	msg.Timestamp = time.Time.UnxNano()
	return msg
}
