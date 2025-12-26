package main

import (
	"fmt"
	"net"
)

func set(conn net.Conn, key, value string) {
	fmt.Println("SET")
	dictionary[key] = value
	encodedResponse := EncodeSimpleString("OK")
	conn.Write(encodedResponse)
}
