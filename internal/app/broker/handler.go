package broker

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kkchu791/sounds/internal/model"
)

type Server struct {
	broker     *model.Broker
	isLeader   bool
	leaderAddr string
}

func NewServer(ID string, p *model.Partition, isLeader bool, leaderAddr string) *Server {
	return &Server{
		broker:     model.NewBroker(ID, p),
		isLeader:   isLeader,
		leaderAddr: leaderAddr,
	}
}

func (s *Server) ReplicateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		offsetStr := req.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)

		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}

		batchLimit := min(offset+10, s.broker.Partition.Len())

		msgs := make([]*model.Message, 0, 10)
		for i := offset; i < batchLimit; i++ {
			msg, err := s.broker.Read(i)
			if err != nil {
				http.Error(w, "had some trouble reading", http.StatusBadRequest)
				return
			}

			msgs = append(msgs, msg)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(msgs)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
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

		message.Timestamp = time.Now()

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
