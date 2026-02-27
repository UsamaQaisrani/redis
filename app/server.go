package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

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
			fmt.Println("Error accepting response = s.Connection: ", err.Error())
			os.Exit(1)
		}
		go server.handleConn()
	}
}

func (s *Server) handleConn() {
	defer s.Conn.Close()

	for {
		command, args := s.readCommand()
		go s.handleCommands(command, args)
	}
}

func (s *Server) readCommand() (command string, args []string) {
	buf := make([]byte, 1024)
	n, err := s.Conn.Read(buf)
	if err != nil {
		return
	}

	fullInput := buf[:n]
	decodedInput, err := Parse(string(fullInput))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	commandArgsAny, ok := decodedInput.([]any)
	if !ok {
		fmt.Println("Unable to convert decoded input to []any")
		return "", nil
	}

	commandArgs := make([]string, len(commandArgsAny))
	for i, v := range commandArgsAny {
		s, ok := v.(string)
		if !ok {
			fmt.Println("Element is not a string:", v)
			return "", nil
		}
		commandArgs[i] = s
	}
	command = commandArgs[0]
	args = commandArgs[1:]
	return command, args
}

func (s *Server) handleCommands(command string, args []string) {
	var response []byte
	switch strings.ToUpper(command) {
	case "PING":
		response = s.Ping("PONG")
	case "ECHO":
		response = s.Echo(args[0])
	case "SET":
		response = s.Set(args)
	case "GET":
		response = s.Get(args[0])
	case "RPUSH":
		response = s.RPush(args)
	case "LRANGE":
		response = s.LRange(args)
	case "LPUSH":
		response = s.LPush(args)
	case "LLEN":
		response = s.LLen(args)
	case "LPOP":
		response = s.LPop(args)
	case "BLPOP":
		response = s.BLPop(args)
	case "TYPE":
		response = s.Type(args)
	case "XADD":
		response = s.XADD(args)
	case "XRANGE":
		streams := s.XRANGE(args)
		response = EncodeStream(streams)
	case "XREAD":
		response = s.XREAD(args)
	case "INCR":
		response = s.Incr(args)
	case "MULTI":
		response = s.Multi(args)
	case "EXEC":
		response = s.Exec()

	default:
		response = []byte("-ERR unknown command\r\n")
	}
	s.Conn.Write(response)
}
