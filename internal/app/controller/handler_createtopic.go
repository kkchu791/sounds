package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/domain/service"
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
	fmt.Println("hey you hit the controller create topic handler")
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

	var results []Result
	for _, topic := range req.Topics {
		mt := model.Topic{
			Name:              topic.Name,
			PartitionCount:    topic.PartitionCount,
			ReplicationFactor: topic.ReplicationFactor,
		}

		err := s.Controller.UpdateTopicPartition(mt)

		var res Result
		if err != nil {
			res = Result{Name: topic.Name, ErrorCode: service.ErrInvalidReplicationFactor, ErrorMsg: err.Error()}
		} else {
			res = Result{Name: topic.Name, ErrorCode: service.ErrNone}
		}

		results = append(results, res)
	}

	res := CreateTopicsResponse{
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
