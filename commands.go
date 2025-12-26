package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func echo(conn net.Conn, arg string) {
	encodedResponse := EncodeBulkString(arg)
	conn.Write(encodedResponse)
}

func set(conn net.Conn, args []string) {
	dictValue := DictStringVal{}
	key := args[0]
	val := args[1]
	if len(args) == 4 {
		timeFormat := args[2]
		exp := args[3]
		createdAt := time.Now().UnixMilli()
		dictValue.CreatedAt = createdAt
		switch strings.ToLower(timeFormat) {
		case "ex":
			seconds, err := strconv.Atoi(exp)
			if err != nil {
				fmt.Println("Unable to convert exp time to string")
				return
			}
			miliseconds := seconds * 1000
			dictValue.Expire = strconv.Itoa(miliseconds)

		case "px":
			dictValue.Expire = exp
		default:
			fmt.Println("Unknow argument for time format")
			return
		}
	}
	dictValue.Value = val
	strDict[key] = dictValue
	encodedResponse := EncodeSimpleString("OK")
	conn.Write(encodedResponse)
}

func get(conn net.Conn, key string) {
	val, ok := strDict[key]
	if !ok {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	if val.Expire != "" {
		createdAt := val.CreatedAt
		expTimer, err := strconv.Atoi(val.Expire)
		if err != nil {
			fmt.Println("Unable to read the value:", err)
			conn.Write([]byte("$-1\r\n"))
			return
		}

		expiryMs := createdAt + int64(expTimer)
		if time.Now().UnixMilli() >= expiryMs {
			delete(strDict, key)
			conn.Write([]byte("$-1\r\n"))
			return
		}
	}

	encodedResponse := EncodeBulkString(val.Value)
	conn.Write(encodedResponse)
}

func ping(conn net.Conn, arg string) {
	encodedResonse := EncodeSimpleString(arg)
	conn.Write(encodedResonse)
}

func rpush(conn net.Conn, args []string) {
	key := args[0]
	vals := args[1:]
	list, ok := listDict[key]
	if !ok {
		newList := []string{}
		list = newList
		listDict[key] = newList
	}
	for _, val := range vals {
		list = append(list, val)
	}
	listDict[key] = list
	encodedResponse := EncodeInt(len(list))
	conn.Write(encodedResponse)
}
