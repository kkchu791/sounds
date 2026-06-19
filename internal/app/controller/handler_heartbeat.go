package controller

import (
	"encoding/json"
	"net/http"
)

type HeartbeatReq struct {
	BrokerID string `json:"broker_id"`
}

func (s *Server) HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// check http method
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var hbr HeartbeatReq
	err := json.NewDecoder(r.Body).Decode(&hbr)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = s.Controller.UpdateLastSeen(hbr.BrokerID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
