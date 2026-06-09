package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kkchu791/sounds/internal/model"
)

type Replicator struct {
	server     *Server
	currOffset int
}

func (r *Replicator) Replicate() {

	// does a GET request on the leader address
	url := fmt.Sprintf("%s/replicate?offset=%d", r.server.leaderAddr, r.currOffset)
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
}
