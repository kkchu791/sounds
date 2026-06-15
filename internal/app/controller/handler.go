package controller

import (
	"encoding/json"
	"log"
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

		s.Controller.Brokers[bi.BrokerID] = bi

		rr := &RegisterResponse{}
		if _, exists := s.Controller.Partitions[partitionID]; !exists {
			s.Controller.Partitions[partitionID] = &model.PartitionInfo{
				LeaderID: bi.BrokerID,
				ISR:      []string{bi.BrokerID},
				Replicas: []string{bi.BrokerID},
			}

			rr.IsLeader = true

			leaderID := s.Controller.Partitions[partitionID].LeaderID
			log.Println(leaderID)
			log.Printf("it ran the leader code in the if block: %+v", rr)

		} else {

			partition := s.Controller.Partitions[partitionID]
			partition.Replicas = append(partition.Replicas, bi.BrokerID)

			rr.IsLeader = false
			leaderID := s.Controller.Partitions[partitionID].LeaderID
			leaderAddr := s.Controller.Brokers[leaderID].BrokerAddr
			rr.LeaderAddr = leaderAddr

			log.Printf("it ran the follower code: %+v", rr)
			log.Printf("leaderID: %s, leaderAddr: %s", leaderID, leaderAddr)
		}

		rr.PartitionID = partitionID

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&rr)
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
