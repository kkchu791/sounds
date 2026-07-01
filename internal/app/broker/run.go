package broker

import (
	"log"
	"net/http"
	"os"

	"github.com/kkchu791/sounds/internal/domain/model"
)

// sets up the server
func Run() {
	id := os.Getenv("BROKER_ID")                   // "broker-0" or "broker-1"
	brokerAddr := os.Getenv("BROKER_ADDR")         // "localhost:5002" or "127.0.0.1:5001"
	port := os.Getenv("BROKER_PORT")               // "5001" or "5002"
	controllerAddr := os.Getenv("CONTROLLER_ADDR") // "localhost:5001" (empty if leader)

	log.Printf("Broker Listening on port: %s", port)

	// this happens first
	resp, err := Register(id, brokerAddr, controllerAddr)

	if err != nil {
		log.Fatalf("failed to register with controller: %s", err)
	}

	server := NewServer(
		id,
		model.NewPartition(),
		resp.IsLeader,
		resp.LeaderAddr,
	)

	// this is for followers to constantly ping the leaders for replication
	if !server.broker.IsLeader {
		r := &Replicator{server: server, currOffset: 0}
		go r.Start()
	}

	// this is for pinging the controller to let it know this broker is alive
	hb := &HeartbeatClient{ID: id, ControllerAddr: controllerAddr}
	go hb.Start()

	http.HandleFunc("/replicate", server.ReplicateHandler)
	http.HandleFunc("/append", server.AppendHandler)
	http.HandleFunc("/read", server.ReadHandler)
	http.HandleFunc("/promote", server.PromoteHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
