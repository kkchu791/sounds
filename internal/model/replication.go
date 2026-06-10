package model

type ReplicationResponse struct {
	HWM      int
	Messages []*Message
}
