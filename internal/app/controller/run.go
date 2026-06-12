package controller

import (
	"log"
	"net/http"
	"time"
)

func Run() {
	server := NewServer()
	mux := http.NewServeMux()

	mux.HandleFunc("/register", server.RegisterHandler)

	httpServer := &http.Server{
		Addr:         ":5000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Server listening on :5000")

	log.Fatal(httpServer.ListenAndServe())
}
