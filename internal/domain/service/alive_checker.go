package service

import (
	"fmt"
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

			bcID, err := ac.HandleDeadBroker(brokerID)
			if err != nil {
				log.Println(err)
				continue
			}

			c := controller.NewClient()
			resp, err := c.Promote(bcID, bcAddr)

			//Cmodel. Update Leader in State if Success
			if resp.ok {
				ac.Controller.UpdateLeader(bcID)
			} else {
				fmt.Printf("failed to promote leader:", bcID)
			}
		}
	}
}
