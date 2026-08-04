package controller

import (
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	BrokerID   string `json:"broker_id"`
	BrokerAddr string `json:"broker_addr"`
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

	err = s.Controller.RegisterBroker(req.BrokerID, req.BrokerAddr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
