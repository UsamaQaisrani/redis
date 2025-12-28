package main

import "sync"

var DB sync.Map

func main() {
	startDB()
	runServer()
}

func startDB() {
	DB = sync.Map{}
}
