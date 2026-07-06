package service

import (
	"context"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/infra/http/broker"
)

type BrokerService struct {
	broker *model.Broker
}

func NewBrokerService(b *model.Broker) *BrokerService {
	return &BrokerService{broker: b}
}

func (s *BrokerService) Register(ctx context.Context, bID, bAddr, cAddr string) error {
	baseURL := "http://" + cAddr
	c := broker.NewClient(baseURL)
	err := c.Register(ctx, bID, bAddr)

	if err != nil {
		return err
	}

	return nil
}

func (s *BrokerService) Promote(isr []string) {
	s.broker.Promote()
	s.broker.UpdateISR(isr)

	// stop the Replcator Stop

	// returns err
}
