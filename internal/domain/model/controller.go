package model

import (
	"sync"
	"time"
)

type Controller struct {
	Brokers    map[string]*BrokerInfo
	Partitions map[string]*PartitionInfo
	mu         sync.Mutex
}

type PartitionInfo struct {
	LeaderID string
	ISR      []string
	Replicas []string
}

type BrokerInfo struct {
	BrokerID   string
	BrokerAddr string
	LastSeen   time.Time
}

type RegisterResult struct {
	IsLeader    bool
	LeaderAddr  string
	PartitionID string
}

func NewController() *Controller {
	return &Controller{
		Brokers:    make(map[string]*BrokerInfo),
		Partitions: make(map[string]*PartitionInfo),
	}
}

func (c *Controller) addBroker(id string, addr string) {
	bi := BrokerInfo{
		BrokerID:   id,
		BrokerAddr: addr,
		LastSeen:   time.Now(),
	}

	c.Brokers[id] = &bi
}

func (c *Controller) addPartition(id string, partitionID string) bool {
	pi, exists := c.Partitions[partitionID]

	if !exists {
		c.Partitions[partitionID] = &PartitionInfo{
			LeaderID: id,
			ISR:      []string{id},
			Replicas: []string{id},
		}

		return true // a leader
	} else {
		pi.Replicas = append(pi.Replicas, id)

		return false // a follower
	}
}

func (c *Controller) RegisterBroker(i string, a string, p string) (*RegisterResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Add Broker
	c.addBroker(i, a)

	//Add to Partition
	isLeader := c.addPartition(i, p)

	var leaderAddr string
	if isLeader {
		leaderAddr = a
	} else {
		leaderID := c.Partitions[p].LeaderID
		leaderAddr = c.Brokers[leaderID].BrokerAddr
	}

	// return something to the handler
	return &RegisterResult{
		IsLeader:    isLeader,
		LeaderAddr:  leaderAddr,
		PartitionID: p,
	}, nil
}
