package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RegisterRequest struct {
	BrokerID   string
	BrokerAddr string
}

type RegisterResponse struct {
	IsLeader    bool
	LeaderAddr  string
	PartitionID string
}

func Register(brokerID, brokerAddr, controllerAddr string) (*RegisterResponse, error) {
	req := RegisterRequest{
		BrokerID:   brokerID,
		BrokerAddr: brokerAddr,
	}

	var buf bytes.Buffer //creates a Buffer object that grows automatically, implements io.Writer

	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://"+controllerAddr+"/register", "application/json", &buf)
	if err != nil {
		log.Printf("post failed: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	var res RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	fmt.Println(&res)

	return &res, nil

}
