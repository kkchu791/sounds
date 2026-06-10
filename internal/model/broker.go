package model

import (
	"log"
	"math"
	"slices"
	"sync"
)

type Broker struct {
	ID              string
	IsLeader        bool
	Partition       *Partition
	HWM             int
	FollowerOffsets map[string]int
	ISR             []string
	mu              sync.Mutex
}

func NewBroker(id string, p *Partition, isLeader bool) *Broker {
	return &Broker{
		ID:              id,
		Partition:       p,
		IsLeader:        isLeader,
		HWM:             0,
		FollowerOffsets: make(map[string]int),
		ISR:             make([]string, 0),
	}
}

func (b *Broker) UpdateFollowerOffset(followerID string, offset int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	currOffset := b.FollowerOffsets[followerID]
	if currOffset >= offset {
		return
	}

	b.FollowerOffsets[followerID] = offset

	bestMin := math.MaxInt

	//only loop througth the ISR and determine best min from them
	for fID, fOffset := range b.FollowerOffsets {
		if slices.Contains(b.ISR, fID) {
			bestMin = min(fOffset, bestMin)
		}
	}

	b.HWM = bestMin - 1
	log.Printf("HWM advanced to %d", b.HWM)

}

func (b *Broker) Len() int {
	return b.Partition.Len()
}

func (b *Broker) Append(m *Message) {
	b.Partition.Append(m) // it just appends the message
}

func (b *Broker) Read(offset int) (*Message, error) {
	return b.Partition.Read(offset)
}
