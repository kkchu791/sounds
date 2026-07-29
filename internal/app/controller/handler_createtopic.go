package controller

import (
	"encoding/json"
	"net/http"
)

// request structs
type CreateTopicsRequest struct {
	Topics []Topic `json:"topics"`
}

type Topic struct {
	Name              string
	PartitionCount    int
	ReplicationFactor int
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
