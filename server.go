package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	Conn net.Conn
}

func runServer() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		server := Server{}
		conn, err := l.Accept()
		server.Conn = conn
		if err != nil {
			fmt.Println("Error accepting s.Connection: ", err.Error())
			os.Exit(1)
		}
		go server.handleConn()
	}
}

func (s Server) handleConn() {
	defer s.Conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := s.Conn.Read(buf)
		if err != nil {
			return
		}

		fullInput := buf[:n]
		decodedInput, err := Parse(string(fullInput))
		if err != nil {
			fmt.Println(err)
			continue
		}

		commandArgsAny, ok := decodedInput.([]any)
		if !ok {
			fmt.Println("Unable to convert decoded input to []any")
			return
		}

		commandArgs := make([]string, len(commandArgsAny))
		for i, v := range commandArgsAny {
			s, ok := v.(string)
			if !ok {
				fmt.Println("Element is not a string:", v)
				continue
			}
			commandArgs[i] = s
		}
		command := commandArgs[0]
		args := commandArgs[1:]

		switch strings.ToLower(command) {
		case "ping":
			s.ping("PONG")
		case "echo":
			s.echo(args[0])
		case "set":
			s.set(args)
		case "get":
			s.get(args[0])
		case "rpush":
			s.rpush(args)
		case "lrange":
			s.lrange(args)
		default:
			s.Conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
