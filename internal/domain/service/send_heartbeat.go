package service

import (
	"context"
	"time"

	"github.com/kkchu791/sounds/internal/infra/http/broker"
)

type SendHeartbeat struct {
	ID             string
	ControllerAddr string
}

func (sh *SendHeartbeat) Start(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		baseURL := "http://" + sh.ControllerAddr
		bClient := broker.NewClient(baseURL)
		bClient.Heartbeat(ctx, sh.ID)
	}
}
