package service

import (
	"log"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
)

type AliveChecker struct {
	Controller *model.Controller
	Timeout    time.Duration
}

func (ac *AliveChecker) Start() {
	c := time.Tick(1 * time.Second)

	for range c {
		deadBrokers := ac.Controller.GetDeadBrokers(ac.Timeout)
		for _, brokerID := range deadBrokers {

			err := ac.Controller.HandleDeadBroker(brokerID)
			if err != nil {
				log.Println(err)
				continue
			}

			//ac.promoteNewleader
		}
	}
}
