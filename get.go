package main

import (
	"fmt"
	"net"
)

func get(conn net.Conn, key string) {
	dictVal, ok := dictionary[key]
	if !ok {
		conn.Write([]byte("$-1\r\n"))
	} else {
		strVal, ok := dictVal.(string)
		if !ok {
			fmt.Println("Unable to convert GET value to string")
			return
		}
		encodedResponse := EncodeBulkString(strVal)
		conn.Write(encodedResponse)
	}
}
