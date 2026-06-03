package model

type Broker struct {
	ID        string
	Partition *Partition
}

func NewBroker(id string, p *Partition) *Broker {
	return &Broker{
		ID:        id,
		Partition: p,
	}
}

func (b *Broker) Append(m *Message) {
	b.Partition.Append(m) // it just appends the message
}

func (b *Broker) Read(offset int) (*Message, error) {
	return b.Partition.Read(offset)
}
