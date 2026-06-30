package service

import (
	"fmt"
	"log"
	"time"

	"github.com/kkchu791/sounds/internal/domain/model"
	"github.com/kkchu791/sounds/internal/infra/http/controller"
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

			bcID, pID, err := ac.Controller.HandleDeadBroker(brokerID)
			if err != nil {
				log.Println(err)
				continue
			}

			bcAddr := ac.Controller.Brokers[bcID].BrokerAddr
			isr := ac.Controller.Partitions[pID].ISR
			c := controller.NewClient(bcAddr)
			pr, err := c.Promote(pID, isr)

			fmt.Print(pr)
			// Cmodel. Update Leader in State if Success
			if pr.Status == "ok" {
				ac.Controller.UpdateLeader(bcID, pID)
			} else {
				fmt.Printf("failed to promote leader: %s", bcID)
			}
		}
	}
}
