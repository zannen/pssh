package main

import (
	"golang.org/x/crypto/ssh"
)

type MultipleConnection struct {
	connections []*Connection
	command     chan string
	response    chan []CmdResponse
}

func NewMultipleConnection() *MultipleConnection {
	mc := MultipleConnection{
		connections: make([]*Connection, 0),
	}

	return &mc
}

func (mc *MultipleConnection) Add(cfg *ssh.ClientConfig, network, hostPort string) {
	c := NewConnection(cfg, network, hostPort)
	mc.connections = append(mc.connections, c)
}

func (mc *MultipleConnection) Command(cmd string) {
	for _, con := range mc.connections {
		con.Command(cmd)
	}
}

func (mc *MultipleConnection) Response() []CmdResponse {
	resp := make([]CmdResponse, len(mc.connections))
	for i, con := range mc.connections {
		resp[i] = con.Response()
	}
	return resp
}
