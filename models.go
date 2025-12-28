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
	Conn   net.Conn
	mu     sync.Mutex
	BLOCKQ []Command
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
