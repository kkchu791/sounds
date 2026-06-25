package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
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

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("failed to create request: %s", err)
		return
	}

	req.Header.Add("X-Broker-ID", r.server.broker.ID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("replication error: %s", err)
		return
	}

	defer resp.Body.Close()

	var rr model.ReplicationResponse
	err = json.NewDecoder(resp.Body).Decode(&rr)
	if err != nil {
		log.Printf("failed to decode replication response: %s", err)
		return
	}

	for _, msg := range rr.Messages {
		r.server.broker.Append(msg)
		r.currOffset++
	}

	// you need update HWM for the follower too here

	if rr.HWM > r.server.broker.HWM {
		r.server.broker.HWM = rr.HWM
	}

	log.Printf("replicated %d messages, now at offset %d, HWM: %d", len(rr.Messages), r.currOffset, r.server.broker.HWM)
}
