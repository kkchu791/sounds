package controller

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/kkchu791/sounds/internal/model"
)

type Server struct {
	Controller *model.Controller
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Controller: model.NewController(),
	}
}

const partitionID = "partition-0"

type RegisterResponse struct {
	IsLeader    bool
	LeaderAddr  string
	PartitionID string
}

func (s *Server) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Method == http.MethodPost {
		var bi *model.BrokerInfo
		err := json.NewDecoder(req.Body).Decode(&bi)

		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		bi.LastSeen = time.Now()

		s.Controller.Brokers[bi.ID] = bi

		rr := &RegisterResponse{}
		if _, exists := s.Controller.Partitions[partitionID]; !exists {
			s.Controller.Partitions[partitionID] = &model.PartitionInfo{
				LeaderID: bi.ID,
				ISR:      []string{bi.ID},
				Replicas: []string{bi.ID},
			}

			rr.IsLeader = true

		} else {

			partition := s.Controller.Partitions[partitionID]
			partition.Replicas = append(partition.Replicas, bi.ID)

			rr.IsLeader = false
			leaderID := s.Controller.Partitions[partitionID].LeaderID
			rr.LeaderAddr = s.Controller.Brokers[leaderID].Address
		}

		rr.PartitionID = partitionID

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&rr)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
