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
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
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

		fmt.Println("DecodedInput:", decodedInput)
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
			ping(conn, "PONG")
		case "echo":
			echo(conn, args[0])
		case "set":
			set(conn, args)
		case "get":
			get(conn, args[0])
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
