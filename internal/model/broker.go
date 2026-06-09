package model

type Broker struct {
	ID        string
	IsLeader  bool
	Partition *Partition
	// HWM             int
	// FollowerOffsets map[string]int
	// ISR             []*Broker
}

func NewBroker(id string, p *Partition, isLeader bool) *Broker {
	return &Broker{
		ID:        id,
		Partition: p,
		IsLeader:  isLeader,
	}
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
