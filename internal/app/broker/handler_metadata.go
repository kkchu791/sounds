package broker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type MetadataResponse struct {
	ControllerAddr string `json:"controller_addr"`
}

func (s *Server) MetadataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ca := os.Getenv("CONTROLLER_ADDR")

	mr := MetadataResponse{
		ControllerAddr: ca,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&mr); err != nil {
		fmt.Println(err)
	}
}
