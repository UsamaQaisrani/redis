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

		commandArgs, ok := decodedInput.([]any)
		if !ok {
			fmt.Println("Unable to convert decode input to list")
			continue
		}

		command, ok := commandArgs[0].(string)
		if !ok {
			fmt.Println("Unable to convert command to string ")
			continue
		}

		switch strings.ToLower(command) {
		case "ping":
			ping(conn, "PONG")
		case "echo":
			arg, ok := commandArgs[1].(string)
			if !ok {
				fmt.Println("Unable to convert arg to string")
				continue
			}
			echo(conn, arg)
		case "set":
			key, ok := commandArgs[1].(string)
			if !ok {
				fmt.Println("Unable to convert SET key to string")
				return
			}
			val, ok := commandArgs[2].(string)
			if !ok {
				fmt.Println("Unable to convert SET value to string")
				return
			}
			set(conn, key, val)
		case "get":
			key, ok := commandArgs[1].(string)
			if !ok {
				fmt.Println("Unable to convert GET key to string")
				return
			}
			get(conn, key)
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
