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
	Conn  net.Conn
	mu    sync.Mutex
	Queue []Command
	DB    DB
}

type Command struct {
	Name     string
	WaitTime int64
}

type DB struct {
	List map[string][]string
	Map  map[string]DictStringVal
}
