package controller

import (
	"time"

	"github.com/kkchu791/sounds/internal/model"
)

func (c *model.Controller) RegisterBroker(bi *model.BrokerInfo, partitionID string) (*RegisterResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bi.LastSeen = time.Now()
	c.Brokers[bi.BrokerID] = bi

	rr := &RegisterResponse{}

	if _, exists := c.Partitions[partitionID]; !exists {
		c.Partitions[partitionID] = &model.PartitionInfo{
			LeaderID: bi.BrokerID,
			ISR:      []string{bi.BrokerID},
			Replicas: []string{bi.BrokerID},
		}

		rr.IsLeader = true
	} else {
		p := c.Partitions[partitionID]
		p.Replicas = append(p.Replicas, bi.BrokerID)

		rr.IsLeader = false

		leaderID := p.LeaderID
		rr.LeaderAddr = c.Brokers[leaderID].BrokerAddr
	}

	rr.PartitionID = partitionID
	return rr, nil
}
