package model

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

const NO_LEADER = ""

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

func (c *Controller) removeBroker(id string) error {
	if _, exists := c.Brokers[id]; exists {
		delete(c.Brokers, id)
	} else {
		return errors.New("can't delete, no broker found")
	}

	return nil
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

func (c *Controller) UpdateLastSeen(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// business logic of updating the state of broker

	_, exists := c.Brokers[id]

	if exists {
		c.Brokers[id].LastSeen = time.Now()
	} else {
		return errors.New("broker not found")
	}

	return nil
}

func (c *Controller) GetDeadBrokers(timeout time.Duration) []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	dead := make([]string, 0)

	for id, val := range c.Brokers {
		// fmt.Printf("this is broker: %v time since last seen value: %v \n", id, time.Since(val.LastSeen))
		// fmt.Printf("this is the configuration timeout we have: %v \n", timeout)
		if time.Since(val.LastSeen) > timeout {
			dead = append(dead, id)
		}
	}

	fmt.Printf("how many dead brokers: %v \n", dead)

	return dead
}

func (c *Controller) HandleDeadBroker(brokerID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.removeBroker(brokerID)
	if err != nil {
		return err
	}

	for id, p := range c.Partitions {
		if idx := slices.Index(p.ISR, brokerID); idx != -1 {
			p.ISR = slices.Delete(p.ISR, idx, idx+1)
		}

		if brokerID == p.LeaderID {
			var newLeader string

			if len(p.ISR) > 0 {
				newLeader = p.ISR[0]
			} else if len(p.Replicas) > 0 {
				for _, r := range p.Replicas {
					if r != brokerID {
						newLeader = r
						break
					}
				}
			} else {
				newLeader = NO_LEADER
			}

			p.LeaderID = newLeader

			fmt.Printf("Leader for partition %s changed from %s to %s\n", id, brokerID, newLeader)
			fmt.Printf("Current ISR for partition %s: %v \n", id, p.ISR)
			fmt.Printf("Current Replicas for partition %s: %v \n", id, p.Replicas)
		} else {
			fmt.Printf("Follower %s removed from partition %s. Current ISR: %s and Replicas: %s \n", brokerID, id, p.ISR, p.Replicas)
		}

	}

	return nil

}
