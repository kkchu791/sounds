package broker

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/domain/service"
)

// sets up the server
func Run() {
	bID := os.Getenv("BROKER_ID")                  // "broker-0" or "broker-1"
	brokerAddr := os.Getenv("BROKER_ADDR")         // "localhost:5002" or "127.0.0.1:5001"
	port := os.Getenv("BROKER_PORT")               // "5001" or "5002"
	controllerAddr := os.Getenv("CONTROLLER_ADDR") // "localhost:5001" (empty if leader)
	ctx := context.Background()

	log.Printf("Broker Listening on port: %s", port)

	var b service.BrokerService
	err := b.Register(ctx, bID, brokerAddr, controllerAddr)
	if err != nil {
		log.Fatalf("failed to register with controller: %s", err)
	}

	const isLeader = false
	const leaderAddr = ""
	broker := model.NewBroker(bID, model.NewPartition(), isLeader)
	mux := http.NewServeMux()
	server := NewServer(
		broker,
		leaderAddr,
		mux,
	)

	// this is for followers to constantly ping the leaders for replication
	// if !server.broker.IsLeader {
	// 	r := &Replicator{server: server, currOffset: 0}
	// 	go r.Start()
	// }

	// this is for pinging the controller to let it know this broker is alive
	hb := service.SendHeartbeat{ID: bID, ControllerAddr: controllerAddr}
	go hb.Start(ctx)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	server.routes()
	log.Fatal(httpServer.ListenAndServe())
}
