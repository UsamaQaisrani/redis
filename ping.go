package main

import "net"

func ping(conn net.Conn, arg string) {
	encodedResonse := EncodeSimpleString(arg)
	conn.Write(encodedResonse)
}
