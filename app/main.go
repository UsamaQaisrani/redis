package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var DB sync.Map

func main() {
	port := flag.Int("port", 6379, "Port to listen on")
	replicaOf := flag.String("replicaof", "", "Address of the master server to replicate from.") 
	flag.Parse()

	if *replicaOf != "" {
		fmt.Printf("Starting server in replica mode, replicating from %s...\n", *replicaOf)
		var parts []string
		if strings.Contains(*replicaOf, ":") {
			parts = strings.Split(*replicaOf, ":")
		} else {
			parts = strings.Split(*replicaOf, " ")
		}
		if len(parts) != 2 {
			fmt.Println("Invalid replicaof address. Expected format: ip:port or ip port")
			return
		}

		_, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			fmt.Println("Invalid port in replicaof address:", err)
			return
		}

		startDB()
		fmt.Printf("Starting slave server on port %d...\n", *port)
		runServer("0.0.0.0", *port, true)
		return
	} else {
		startDB()
		fmt.Printf("Starting server on port %d...\n", *port)
		runServer("0.0.0.0", *port, false)
	} 

}

func startDB() {
	DB = sync.Map{}
}
