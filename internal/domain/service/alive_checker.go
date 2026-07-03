package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/infra/http/controller"
)

type AliveChecker struct {
	Controller *model.Controller
	Timeout    time.Duration
}

func (ac *AliveChecker) Start(ctx context.Context) { // it runs alive checker
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ac.HandleDeadBrokers(ctx)
	}
}

func (ac *AliveChecker) HandleDeadBrokers(ctx context.Context) {
	list := ac.Controller.GetDeadBrokers(ac.Timeout)

	for _, dead := range list {
		bcID, bcAddr, pID, isr := ac.Controller.HandleDeadBroker(dead)

		if bcID != "" {
			ac.Promote(ctx, bcID, bcAddr, pID, isr)
		}
	}
}

func (ac *AliveChecker) Promote(ctx context.Context, bcID, bcAddr, pID string, isr []string) {
	baseURL := "http://" + bcAddr
	c := controller.NewClient(baseURL)
	pr, err := c.Promote(ctx, pID, isr)

	if err != nil {
		fmt.Printf("promoting error: %v", err)
		return
	}

	if pr.Status == "ok" {
		ac.Controller.UpdateLeader(bcID, pID)
	} else {
		fmt.Printf("failed to promote leader: %s", bcID)
	}
}
