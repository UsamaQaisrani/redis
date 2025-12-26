package main

import "net"

func echo(conn net.Conn, arg string) {
	encodedResponse := EncodeBulkString(arg)
	conn.Write(encodedResponse)
}
