package broker

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/kkchu791/sounds/internal/model"
)

// sets up the server
func Run() {

	id := os.Getenv("BROKER_ID")                               // "broker-0" or "broker-1"
	isLeader, err := strconv.ParseBool(os.Getenv("IS_LEADER")) // "true" or "false"
	port := os.Getenv("BROKER_PORT")                           // "5001" or "5002"
	leaderAddr := os.Getenv("LEADER_ADDR")                     // "localhost:5001" (empty if leader)

	if err != nil {
		log.Fatalf("invalid IS_LEADER value: %s", err)
	}

	log.Printf("Broker Listening on port: %s", port)

	// create a customized broker based on env var
	server := NewServer(
		id,
		model.NewPartition(),
		isLeader,
		leaderAddr,
	)

	// this is for followers to constantly ping the leaders for replication
	if !server.broker.IsLeader {
		r := &Replicator{server: server, currOffset: 0}
		go r.Start()
	}

	http.HandleFunc("/replicate", server.ReplicateHandler)
	http.HandleFunc("/append", server.AppendHandler)
	http.HandleFunc("/read", server.ReadHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
