package main

import (
	"net"
	"sync"
)

type DictStringVal struct {
	Value     string
	Expire    string
	CreatedAt int64
}

type Server struct {
	Conn    net.Conn
	mu      sync.Mutex
	BLOCKQ  []Command
	TxQueue [][]string
	SType ServerType
}

type Data struct {
	Content   any
	ExpiresAt int64
	Waiting   map[string][]chan string
}

type Command struct {
	Name     string
	WaitTime int64
}

type Stream struct {
	StreamID      string
	KeyValuePairs map[string]string
}

type ServerType struct {
	Role string
	ip  string
	port int
}
