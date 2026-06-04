package broker

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/kkchu791/sounds/internal/model"
)

type Server struct {
	broker *model.Broker
}

func NewServer(ID string, p *model.Partition) *Server {
	return &Server{
		broker: model.NewBroker(ID, p),
	}
}

func (s *Server) ReadHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		offsetStr := req.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)

		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		msg, err := s.broker.Read(offset)

		if err != nil {
			http.Error(w, "had some trouble reading", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(msg)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// individual handler func fo each route
func (s *Server) AppendHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var message model.Message
		err := json.NewDecoder(req.Body).Decode(&message)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		s.broker.Append(&message)

		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "message appended successfully")
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
