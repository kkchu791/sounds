package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kkchu791/sounds/internal/domain/model"
)

// request structs
type CreateTopicsRequest struct {
	Topics []Topic `json:"topics"`
}

type Topic struct {
	Name              string `json:"topic"`
	PartitionCount    int    `json:"partition_count"`
	ReplicationFactor int    `json:"replication_factor"`
}

// response structs
type CreateTopicsResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name      string `json:"name"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func (s *Server) CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTopicsRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// TODO: add logic here
	// takes in a topic configurations like topic name partition count, replication factor
	//result, err := s.Controller.CreateTopic()
	// Topics []Topic

	for _, topic := range req.Topics {

		// validate if rf is smaller than servers

		if topic.ReplicationFactor > len(s.Controller.Brokers) {
			fmt.Println("rf too large, need to add brokers or decrease rf")
		}

		s.Controller.UpdateTopicPartition(model.Topic{
			Name:              topic.Name,
			PartitionCount:    topic.PartitionCount,
			ReplicationFactor: topic.ReplicationFactor,
		})

	}

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//just dummy data for now
	result := Result{Name: "Sounds", ErrorCode: 0}

	results := []Result{result}

	res := CreateTopicsResponse{
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
