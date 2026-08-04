package controller

import (
	"log"
	"net/http"
	"time"
)

func Run() {
	server := NewServer()
	mux := http.NewServeMux()

	//TODO: will bring back for shutting off a node
	// ctx := context.Background()

	// this is a worker for checking if any broker nodes are dead and handle leader election
	// ac := &service.AliveChecker{
	// 	Controller: server.Controller,
	// 	Timeout:    time.Duration(timeout) * time.Second,
	// }
	// go ac.Start(ctx)

	mux.HandleFunc("/register", server.RegisterHandler)
	mux.HandleFunc("/heartbeat", server.HeartbeatHandler)
	mux.HandleFunc("/createTopic", server.CreateTopicHandler)

	httpServer := &http.Server{
		Addr:         ":9000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Server listening on :9000")

	log.Fatal(httpServer.ListenAndServe())
}
