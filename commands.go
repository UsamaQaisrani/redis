package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s *Server) Echo(arg string) []byte {
	encodedResponse := EncodeBulkString(arg)
	return encodedResponse
}

func (s *Server) Set(args []string) []byte {
	data := Data{}
	key := args[0]
	val := args[1]
	if len(args) == 4 {
		timeFormat := args[2]
		exp := args[3]
		ms, err := strconv.Atoi(exp)
		if err != nil {
			return nil
		}
		if strings.ToUpper(timeFormat) == "PX" {
			data.ExpiresAt = time.Now().Add(time.Duration(ms) * time.Millisecond).UnixMilli()
		}
	}
	data.Content = val
	DB.Store(key, data)
	encodedResponse := EncodeSimpleString("OK")
	return encodedResponse
}
func (s *Server) Get(key string) []byte {
	val, ok := DB.Load(key)
	if !ok {
		return []byte("$-1\r\n")
	}

	data := val.(Data)

	if data.ExpiresAt != 0 {
		fmt.Println("Has Expiry:", data.ExpiresAt)
		if time.Now().UnixMilli() >= data.ExpiresAt {
			fmt.Println("Expired")
			DB.Delete(key)
			return []byte("$-1\r\n")
		}
		fmt.Println("Not Expired")
	}

	var encodedResponse []byte

	switch data.Content.(type) {
	case string:
		strRes := data.Content.(string)
		encodedResponse = EncodeBulkString(strRes)
	default:
		fmt.Println("Unkown type for the value")
	}

	return encodedResponse
}

func (s *Server) Ping(arg string) []byte {
	encodedResonse := EncodeSimpleString(arg)
	return encodedResonse
}

func (s *Server) RPush(args []string) []byte {
	key := args[0]
	vals := args[1:]
	dbVal, ok := DB.Load(key)
	if !ok {
		dbVal = Data{Content: []string{}}
	}
	data := dbVal.(Data)
	list := data.Content.([]string)
	for _, val := range vals {
		list = append(list, val)
	}

	chs := data.Waiting["RPUSH"]
	if len(chs) > 0 {
		ch := chs[0]
		data.Waiting["RPUSH"] = chs[1:]
		ch <- vals[0]
	}

	data.Content = list
	DB.Store(key, data)
	encodedResponse := EncodeInt(len(list))
	return encodedResponse
}

func (s *Server) LRange(args []string) []byte {
	key := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Unable to parse the range start")
		return nil
	}
	end, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Unable to parse the range end")
		return nil
	}

	dbVal, ok := DB.Load(key)
	if !ok {
		encodedResponse := EncodeList([]string{})
		return encodedResponse
	}
	data := dbVal.(Data)
	list := data.Content.([]string)

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
		return encodedResponse
	}

	if end >= len(list) {
		end = len(list) - 1
	}

	encodedResponse := EncodeList(list[start : end+1])
	return encodedResponse
}

func (s *Server) LPush(args []string) []byte {
	key := args[0]
	items := args[1:]
	dbVal, ok := DB.Load(key)
	if !ok {
		dbVal = Data{Content: []string{}}
	}
	data := dbVal.(Data)
	list := data.Content.([]string)
	for _, item := range items {
		list = append([]string{item}, list...)
	}
	data.Content = list
	DB.Store(key, data)
	encodedResponse := EncodeInt(len(list))
	return encodedResponse
}

func (s *Server) LLen(args []string) []byte {
	key := args[0]
	dbVal, ok := DB.Load(key)
	if !ok {
		dbVal = Data{Content: []string{}}
	}
	data := dbVal.(Data)
	list := data.Content.([]string)
	encodedResponse := EncodeInt(len(list))
	return encodedResponse
}

func (s *Server) LPop(args []string) []byte {
	key := args[0]
	itemsToRemove := 1
	if len(args) == 2 {
		itemsToRemove, _ = strconv.Atoi(args[1])
	}
	dbVal, ok := DB.Load(key)
	if !ok {
		return []byte("$-1\r\n")
	}
	data := dbVal.(Data)
	list := data.Content.([]string)
	if len(list) < 1 {
		return []byte("$-1\r\n")
	}

	poppedItem := list[:itemsToRemove]
	data.Content = list[itemsToRemove:]
	DB.Store(key, data)

	var encodedResponse []byte

	if len(poppedItem) > 1 {
		encodedResponse = EncodeList(poppedItem)
	} else {
		encodedResponse = EncodeBulkString(poppedItem[0])
	}

	return encodedResponse
}

func (s *Server) BLPop(args []string) []byte {
	fmt.Println(args)
	key := args[0]
	wait, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return []byte("Unable to parse the time interval to int")
	}
	dbVal, ok := DB.Load(key)
	var data Data
	if !ok {
		data = Data{
			Content: make([]string, 0),
			Waiting: make(map[string][]chan string),
		}
	} else {
		data = dbVal.(Data)
		if data.Waiting == nil {
			data.Waiting = make(map[string][]chan string)
		}
	}

	var timeOutCh <-chan time.Time
	ch := make(chan string)
	data.Waiting["RPUSH"] = append(data.Waiting["RPUSH"], ch)
	DB.Store(key, data)

	if wait > 0 {
		timeOutCh = time.After(time.Duration(wait * float64(time.Second)))
	}

	var encodedResponse []byte
	select {
	case item := <-ch:
		s.LPop([]string{key})
		respList := []string{key, item}
		encodedResponse = EncodeList(respList)
	case <-timeOutCh:
		encodedResponse = []byte("*-1\r\n")
	}
	return encodedResponse
}
func (s *Server) Type(args []string) []byte {
	key := args[0]
	dbVal, ok := DB.Load(key)
	var encodedResponse []byte
	if !ok {
		fmt.Printf("Content Type: %T", dbVal)
		return EncodeSimpleString("none")
	}

	data := dbVal.(Data)

	fmt.Printf("Content Type: %T", data.Content)

	switch data.Content.(type) {
	case string:
		encodedResponse = EncodeSimpleString("string")
	case []string:
		encodedResponse = EncodeSimpleString("list")
	case map[string][]Stream:
		encodedResponse = EncodeSimpleString("stream")
	}
	return encodedResponse
}

func (s *Server) XADD(args []string) []byte {
	var encodedResponse []byte
	key := args[0]
	id := args[1]
	dbVal, ok := DB.Load(key)
	if !ok {
		dbVal = Data{Content: map[string][]Stream{}}
	}

	if id == "0-0" {
		encodedResponse = EncodeSimpleError("ERR The ID specified in XADD must be greater than 0-0")
		return encodedResponse
	}

	data := dbVal.(Data)
	streams := data.Content.(map[string][]Stream)

	newId, err := generateStreamId(streams, key, id)
	if err != nil && newId == "" {
		encodedResponse = EncodeSimpleError("ERR The ID specified in XADD is equal or smaller than the target stream top item")
		return encodedResponse
	}

	pairs := map[string]string{}
	i := 2

	for i < len(args) {
		k := args[i]
		v := args[i+1]
		pairs[k] = v
		i += 2
	}

	streams[key] = append(streams[key], Stream{StreamID: newId, KeyValuePairs: pairs})
	data.Content = streams
	DB.Store(key, data)
	return EncodeBulkString(newId)
}

func (s *Server) XRANGE(args []string) []byte {
	key := args[0]

	// The command can accept IDs in the format <millisecondsTime>-<sequenceNumber>,
	// but the sequence number is optional.
	var start string
	end := args[2]

	if strings.Contains(args[1], "-") {
		if args[1] == "-" {
			// Start from the beginning of the stream
			start = "0-0"
		} else {
			// Start from the id sent by the user
			start = args[1]
		}

	} else {
		start = args[1] + "-0"
	}

	dbVal, ok := DB.Load(key)
	if !ok {
		return nil
	}

	data := dbVal.(Data)
	streams, ok := data.Content.(map[string][]Stream)
	if !ok {
		return nil
	}

	stream := streams[key]

	if !strings.Contains(end, "-") {
		end += "-" + strconv.Itoa(len(stream))
	}

	var streamSlice []Stream
	for _, entry := range stream {
		if inRangeStreamId(entry.StreamID, start, end) {
			streamSlice = append(streamSlice, entry)
		}
	}

	return EncodeStream(streamSlice)
}
