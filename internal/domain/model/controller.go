package model

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kkchu791/sounds/internal/domain/utils"
)

const NO_LEADER = ""

type Controller struct {
	Brokers    map[string]*BrokerInfo
	Topics     map[string]*TopicInfo
	Partitions map[PartitionName]*PartitionInfo
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

type PartitionName struct {
	TopicName   string
	PartitionId int
}

type PartitionInfo struct {
	LeaderID       int
	ISR            []int
	Replicas       []int
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
		Partitions: make(map[PartitionName]*PartitionInfo),
	}
}

func (c *Controller) RegisterBroker(id string, addr string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Add Broker

	_, exists := c.Brokers[id]
	if exists {
		return fmt.Errorf("broker already exists")
	} else {
		bi := BrokerInfo{
			BrokerID:   id,
			BrokerAddr: addr,
			LastSeen:   time.Now(),
		}

		c.Brokers[id] = &bi
	}

	return nil
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

func (c *Controller) UpdateTopicPartition(topic Topic) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	pc := topic.PartitionCount
	tn := topic.Name
	rf := topic.ReplicationFactor
	bn := len(c.Brokers)

	if rf > len(c.Brokers) {
		return fmt.Errorf("rf too large, need to add brokers or decrease rf")
	}

	c.Topics[topic.Name] = &TopicInfo{
		Name:              tn,
		PartitionCount:    pc,
		ReplicationFactor: rf,
	}

	// Do RR + Shift
	mapPToR := utils.AssignReplicasToBrokers(pc, rf, bn, -1, 1)

	// create the Partition
	for pId, r := range mapPToR {
		key := PartitionName{TopicName: tn, PartitionId: pId}
		value := PartitionInfo{
			LeaderID:       r[0],
			Replicas:       r,
			ISR:            r,
			LeaderEpoch:    0,
			PartitionEpoch: 0,
		}

		c.Partitions[key] = &value
	}
	fmt.Println("-----")

	for k, v := range c.Partitions {
		fmt.Println(k)
		fmt.Println(*v)
	}

	return nil
}
