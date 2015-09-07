package dht

import {
	"fmt"
	"encoding/json"
}

type Transport struct{
	bindAddress string 
}

func (transport	*Transport)	listen(){
	udpAddr, err :=	net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec	:=	json.NewDecoder(conn)
	for	{
		msg :=	Msg{}
		err	=	dec.Decode(&msg)
//	we	got	a	message, do something
	}
}

func (transport *Transport) send(msg *Msg) {
	udpAddr,	err	:=	net.ResolveUDPAddr("udp",	dhtMsg.Dst)
	conn,	err	:=	net.DialUDP("udp",	nil,	udpAddr)
	defer	conn.Close()
	_,	err	=	conn.Write(msg.Bytes())
}