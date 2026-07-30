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
	Topics     map[string]*TopicInfo
	Partitions map[string]*PartitionInfo
	mu         sync.Mutex
}

type BrokerInfo struct {
	BrokerID   string
	BrokerAddr string
	LastSeen   time.Time
	Fenced     bool
}

type TopicInfo struct {
	Name              string
	PartitionCount    int
	ReplicationFactor int
}
type Topic struct {
	Name              string
	PartitionCount    int
	ReplicationFactor int
}

type PartitionInfo struct {
	Topic          string
	PartitionIndex int
	LeaderID       string
	ISR            []string
	Replicas       []string
	LeaderEpoch    int
	PartitionEpoch int
}

type RegisterResult struct {
	IsLeader    bool
	LeaderAddr  string
	PartitionID string
}

func NewController() *Controller {
	return &Controller{
		Brokers:    make(map[string]*BrokerInfo),
		Topics:     make(map[string]*TopicInfo),
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

func (c *Controller) addPartition(bID string, partitionID string) bool {
	pi, exists := c.Partitions[partitionID]

	if !exists {
		c.Partitions[partitionID] = &PartitionInfo{
			LeaderID: bID,
			ISR:      []string{bID},
			Replicas: []string{bID},
		}

		return true // a leader
	} else {
		if !slices.Contains(pi.Replicas, bID) {
			pi.Replicas = append(pi.Replicas, bID)
		}

		return false // a follower
	}
}

func (c *Controller) RegisterBroker(i string, a string, p string) (*RegisterResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Add Broker
	c.addBroker(i, a)

	// this usually occurs when adminclient creates a topic.
	// kafka-topics.sh --create --topic sounds --partitions 3
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

	deadBrokers := make([]string, 0)

	for id, val := range c.Brokers {
		// fmt.Printf("this is broker: %v time since last seen value: %v \n", id, time.Since(val.LastSeen))
		// fmt.Printf("this is the configuration timeout we have: %v \n", timeout)
		if time.Since(val.LastSeen) > timeout {
			deadBrokers = append(deadBrokers, id)
		}
	}

	fmt.Printf("dead brokers: %v \n", deadBrokers)

	return deadBrokers
}

func (c *Controller) HandleDeadBroker(brokerID string) (string, string, string, []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeBroker(brokerID)

	// TODO: only handles single partition for now, need to refactor if more partitions are added
	var bcID string
	var pID string
	var isr []string
	var bcAddr string
	for id, p := range c.Partitions {
		c.removeBrokerFromISRAndReplica(p, brokerID)

		if brokerID == p.LeaderID {
			bcID = c.selectBestCandidate(p, brokerID)

			fmt.Printf("Potential Leader for partition %s changed from %s to %s \n", id, brokerID, bcID)
			fmt.Printf("Current ISR for partition %s: %v \n", id, p.ISR)
			fmt.Printf("Current Replicas for partition %s: %v \n", id, p.Replicas)
		} else {
			fmt.Printf("Follower %s removed from partition %s. Current ISR: %s and Replicas: %s \n", brokerID, id, p.ISR, p.Replicas)
		}

		pID = id
		isr = p.ISR
	}

	if bcID != "" {
		bcAddr = c.Brokers[bcID].BrokerAddr
	}

	return bcID, bcAddr, pID, isr

}

func (c *Controller) removeBrokerFromISRAndReplica(p *PartitionInfo, bID string) {
	if idx := slices.Index(p.ISR, bID); idx != -1 {
		p.ISR = slices.Delete(p.ISR, idx, idx+1)
	}

	// if idx := slices.Index(p.Replicas, bID); idx != -1 {
	// 	p.Replicas = slices.Delete(p.Replicas, idx, idx+1)
	// }
}

func (c *Controller) selectBestCandidate(p *PartitionInfo, bID string) string {
	var bcID string
	if len(p.ISR) > 0 {
		bcID = p.ISR[0]
	} else if len(p.Replicas) > 0 {
		for _, r := range p.Replicas {
			if r != bID && c.isBrokerAlive(r) {
				bcID = r
				break
			}
		}
	}

	if bcID == "" {
		fmt.Println("found no eligible leaders to promote")
	}

	return bcID
}

func (c *Controller) isBrokerAlive(bID string) bool {
	_, exists := c.Brokers[bID]
	return exists
}

func (c *Controller) UpdateLeader(bcID, pID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	pi := c.Partitions[pID]
	pi.LeaderID = bcID

	fmt.Printf("partition %s's leader has been updated to: %s \n", pID, bcID)
	fmt.Println("promotion completed, epoch ended")
}

func (c *Controller) UpdateTopicPartition(topic Topic) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Topics[topic.Name] = &TopicInfo{
		Name:              topic.Name,
		PartitionCount:    topic.PartitionCount,
		ReplicationFactor: topic.ReplicationFactor,
	}
}
