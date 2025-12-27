package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s Server) echo(arg string) {
	encodedResponse := EncodeBulkString(arg)
	s.Conn.Write(encodedResponse)
}

func (s Server) set(args []string) {
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
	s.Conn.Write(encodedResponse)
}

func (s Server) get(key string) {
	val, ok := strDict[key]
	if !ok {
		s.Conn.Write([]byte("$-1\r\n"))
		return
	}

	if val.Expire != "" {
		createdAt := val.CreatedAt
		expTimer, err := strconv.Atoi(val.Expire)
		if err != nil {
			fmt.Println("Unable to read the value:", err)
			s.Conn.Write([]byte("$-1\r\n"))
			return
		}

		expiryMs := createdAt + int64(expTimer)
		if time.Now().UnixMilli() >= expiryMs {
			delete(strDict, key)
			s.Conn.Write([]byte("$-1\r\n"))
			return
		}
	}

	encodedResponse := EncodeBulkString(val.Value)
	s.Conn.Write(encodedResponse)
}

func (s Server) ping(arg string) {
	encodedResonse := EncodeSimpleString(arg)
	s.Conn.Write(encodedResonse)
}

func (s Server) rpush(args []string) {
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
	s.Conn.Write(encodedResponse)
}

func (s Server) lrange(args []string) {
	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Unable to parse the range start")
		return
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Unable to parse the range end")
		return
	}

	list, ok := listDict[key]
	if !ok {
		list = []string{}
		encodedResponse := EncodeList([]string{})
		s.Conn.Write(encodedResponse)
		return
	}

	if start < 0 {
		adjustedStart := len(list) + start
		start = max(adjustedStart, 0)
	}

	if end < 0 {
		end = len(list) + end
	}

	if start >= len(list) || start > end {
		list = []string{}
		encodedResponse := EncodeList([]string{})
		s.Conn.Write(encodedResponse)
		return
	}

	if end >= len(list) {
		end = len(list) - 1
	}

	encodedResponse := EncodeList(list[start : end+1])
	s.Conn.Write(encodedResponse)
}
