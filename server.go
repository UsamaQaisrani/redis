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
			s.Ping("PONG")
		case "echo":
			s.Echo(args[0])
		case "set":
			s.Set(args)
		case "get":
			s.Get(args[0])
		case "rpush":
			s.RPush(args)
		case "lrange":
			s.LRange(args)
		case "lpush":
			s.LPush(args)
		case "llen":
			s.LLen(args)
		case "lpop":
			s.LPop(args)
		default:
			s.Conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
