package model

import "time"

type Message struct {
	Sound     string
	Key       string
	Timestamp time.Time
}

func NewMessage(m, k string) *Message {
	return &Message{
		Sound:     m,
		Key:       k,
		Timestamp: time.Now(),
	}
}
