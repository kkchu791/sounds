package broker

import (
	"log"
	"net/http"

	"github.com/kkchu791/sounds/internal/model"
)

// sets up the server
func Run() {
	log.Println("Broker Listening on port 5001")

	// create a new server with a new partition
	server := NewServer("broker-0", model.NewPartition())

	http.HandleFunc("/append", server.AppendHandler)
	http.HandleFunc("/read", server.ReadHandler)
	log.Fatal(http.ListenAndServe(":5001", nil))

}
