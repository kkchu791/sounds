package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkchu791/sounds/internal/model"
)

type Replicator struct {
	server     *Server
	currOffset int
}

func (r *Replicator) Start() {
	c := time.Tick(1 * time.Second) // returns a channel with time sent into it
	for range c {
		r.Replicate()
	}
}

func (r *Replicator) Replicate() {

	// does a GET request on the leader address
	url := fmt.Sprintf("http://%s/replicate?offset=%d", r.server.leaderAddr, r.currOffset)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("replication error: %s", err)
		return
	}

	defer resp.Body.Close()

	var messages []*model.Message
	json.NewDecoder(resp.Body).Decode(&messages)

	for _, msg := range messages {
		r.server.broker.Append(msg)
		r.currOffset++
	}

	log.Printf("replicated %d messages, now at offset %d", len(messages), r.currOffset)
}
