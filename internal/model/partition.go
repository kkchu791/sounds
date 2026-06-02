package model

import "fmt"

type Partition struct {
	messages []*Message
}

func NewPartition() *Partition {
	return &Partition{
		messages: make([]*Message, 0),
	}
}

func (p *Partition) Append(m *Message) {
	p.messages = append(p.messages, m)
}

func (p *Partition) Read(offset int) (*Message, error) {
	if offset >= len(p.messages) || offset < 0 {
		err := fmt.Errorf("hey this offset is out of range")
		return nil, err
	}

	return p.messages[offset], nil
}
