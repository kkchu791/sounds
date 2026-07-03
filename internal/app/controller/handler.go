package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kkchu791/sounds/internal/domain/model"
)

type Server struct {
	Controller *model.Controller
}

func NewServer() *Server {
	return &Server{
		Controller: model.NewController(),
	}
}

type RegisterRequest struct {
	BrokerID   string `json:"broker_id"`
	BrokerAddr string `json:"broker_addr"`
}
type RegisterResponse struct {
	IsLeader    bool   `json:"is_leader"`
	LeaderAddr  string `json:"leader_addr"`
	PartitionID string `json:"partition_id"`
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := s.Controller.RegisterBroker(req.BrokerID, req.BrokerAddr, partitionID)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	res := RegisterResponse{
		IsLeader:    result.IsLeader,
		LeaderAddr:  result.LeaderAddr,
		PartitionID: result.PartitionID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
