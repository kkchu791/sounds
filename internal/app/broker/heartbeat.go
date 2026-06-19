package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HeartbeatClient struct {
	ID             string
	ControllerAddr string
}

type HeartbeatReq struct {
	BrokerID string `json:"broker_id"`
}

func (hb *HeartbeatClient) Start() {
	c := time.Tick(3 * time.Second)
	for range c {
		err := Heartbeat(hb.ID, hb.ControllerAddr)
		if err != nil {
			log.Println(err)
		}

		// TODO: Lets handle the retry
		// if resp != "ok" {
		// 	fmt.Println("I'm aiming to retry the register if it fails")
		// }
	}
}

func Heartbeat(id string, cAddr string) error {
	req := HeartbeatReq{
		BrokerID: id,
	}

	var buf bytes.Buffer //creates a Buffer object that grows automatically, implements io.Writer

	e := json.NewEncoder(&buf)
	err := e.Encode(&req)

	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+cAddr+"/heartbeat", "application/json", &buf)
	if err != nil {
		log.Printf("post failed: %v", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status: %d", resp.StatusCode)
	}

	return nil
}
