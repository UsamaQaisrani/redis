package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func runServer() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		server := Server{}
		server.DB.List = map[string][]string{}
		server.DB.Map = map[string]DictStringVal{}
		conn, err := l.Accept()
		server.Conn = conn
		if err != nil {
			fmt.Println("Error accepting s.Connection: ", err.Error())
			os.Exit(1)
		}
		go server.handleConn()
	}
}

func (s *Server) handleConn() {
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
		go s.handleCommands(command, args)
	}
}

func (s *Server) handleCommands(command string, args []string) {
	switch strings.ToLower(command) {
	case "ping":
		cmd := Command{Name: "PING", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.Ping("PONG")
	case "echo":
		cmd := Command{Name: "ECHO", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.Echo(args[0])
	case "set":
		cmd := Command{Name: "SET", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.Set(args)
	case "get":
		cmd := Command{Name: "GET", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.Get(args[0])
	case "rpush":
		cmd := Command{Name: "RPUSH", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.RPush(args)
	case "lrange":
		cmd := Command{Name: "LRANGE", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.LRange(args)
	case "lpush":
		cmd := Command{Name: "LPUSH", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.LPush(args)
	case "llen":
		cmd := Command{Name: "LLEN", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.LLen(args)
	case "lpop":
		cmd := Command{Name: "LPOP", WaitTime: int64(0)}
		s.Queue = append(s.Queue, cmd)
		s.LPop(args)
	case "blpop":
		waitingTimeSec, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Unable to convert time for BLPOP")
			return
		}
		waitingTimeMil := int64(waitingTimeSec * 1000)
		currTime := time.Now().UnixMilli()
		cmd := Command{Name: "BLPOP", WaitTime: currTime + waitingTimeMil}
		s.Queue = append(s.Queue, cmd)
		s.BLPop(args)
	default:
		s.Conn.Write([]byte("-ERR unknown command\r\n"))
	}
}
