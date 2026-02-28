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

	startDB()

	if *replicaOf != "" {
		runReplica(*port, *replicaOf)
	} else {
		runMaster(*port)
	}
}

func runMaster(port int) {
	fmt.Printf("Starting server on port %d...\n", port)
	runServer("0.0.0.0", port, false)
}

func runReplica(port int, replicaOf string) {
	cfg, err := parseReplicaOf(replicaOf)
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.SelfPort = strconv.Itoa(port)
	fmt.Printf("Starting server in replica mode, replicating from %s...\n", cfg.MasterAddr())
	masterConn, err := Handshake(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer masterConn.Close()
	fmt.Printf("Starting slave server on port %d...\n", port)
	runServer("0.0.0.0", port, true)
}

func parseReplicaOf(replicaOf string) (*ReplicationConfig, error) {
	var parts []string
	if strings.Contains(replicaOf, ":") {
		parts = strings.Split(replicaOf, ":")
	} else {
		parts = strings.Split(replicaOf, " ")
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid replicaof address. Expected format: ip:port or ip port")
	}
	host := strings.TrimSpace(parts[0])
	port := strings.TrimSpace(parts[1])
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("invalid port in replicaof address: %w", err)
	}
	return &ReplicationConfig{MasterHost: host, MasterPort: port}, nil
}

func startDB() {
	DB = sync.Map{}
}
