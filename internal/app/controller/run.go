package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/kkchu791/sounds/internal/domain/service"
)

func Run() {
	server := NewServer()
	mux := http.NewServeMux()

	// this is a worker for checking if any broker nodes are dead and handle leader election
	ac := &service.AliveChecker{
		Controller: server.Controller,
		Timeout:    time.Duration(timeout) * time.Second,
	}
	go ac.Start()

	mux.HandleFunc("/register", server.RegisterHandler)
	mux.HandleFunc("/heartbeat", server.HeartbeatHandler)

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
