package broker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/domain/service"
)

type Server struct {
	broker     *model.Broker
	leaderAddr string
	mux        *http.ServeMux
}

func NewServer(b *model.Broker, lAddr string, m *http.ServeMux) *Server {
	return &Server{
		broker:     b,
		leaderAddr: lAddr,
		mux:        m,
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

		followerID := req.Header.Get("X-Broker-ID")

		s.broker.UpdateFollowerOffset(followerID, offset)

		// get data from partition
		batchLimit := 10
		batchOffset := min(offset+batchLimit, s.broker.Partition.Len())

		msgs := make([]*model.Message, 0, batchLimit)
		for i := offset; i < batchOffset; i++ {
			msg, err := s.broker.Read(i)
			if err != nil {
				http.Error(w, "had some trouble reading", http.StatusBadRequest)
				return
			}

			msgs = append(msgs, msg)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		rr := model.ReplicationResponse{
			HWM:      s.broker.HWM,
			Messages: msgs,
		}

		json.NewEncoder(w).Encode(rr)
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

type PromoteRequest struct {
	PartitionID string   `json:"partition_id"`
	ISR         []string `json:"isr"`
}
type PromoteResponse struct {
	Status string `json:"status"`
}

func (s *Server) PromoteHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var pr PromoteRequest
		err := json.NewDecoder(req.Body).Decode(&pr)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Broker Serve promote
		service := service.NewBrokerService(s.broker)
		service.Promote(pr.ISR)

		fmt.Println(s.broker)

		//write status ok
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		pres := PromoteResponse{
			Status: "ok",
		}
		json.NewEncoder(w).Encode(&pres)

	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
