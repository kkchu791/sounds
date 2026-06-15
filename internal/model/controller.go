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

func NewController() *Controller {
	return &Controller{
		Brokers:    make(map[string]*BrokerInfo),
		Partitions: make(map[string]*PartitionInfo),
	}
}
