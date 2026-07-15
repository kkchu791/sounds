package service

import (
	"fmt"

	"github.com/kkchu791/sounds/internal/infra/http/controller"
)

type CreateTopicsResponse struct {
	Topics []CreateTopicResult
}

const (
	ErrNone                     = 0
	ErrInvalidReplicationFactor = 38
)

type CreateTopicResult struct {
	Name      string
	ErrorCode int
	ErrorMsg  string
}

func CreatTopic(tn, pc, rf, bAddr string) (CreateTopicsResponse, error) {
	fmt.Println("in the service, Hey")

	//creates a client for broker
	aClient := admin.NewClient(bAddr)
	resp, err := aClient.GetMetaData()

	if err != nil {
		fmt.Println(err)
		return err
	}

	//resp should equal MetaDataRes{controller_addr: "localhost:5000"}

	cAddr := resp.controller_addr
	cClient := controller.NewClient(cAddr)

	resp, err = cClient.CreateTopic()

	if err != nil {
		fmt.Println(err)
		return err
	}
	// resp should equal
	// CreateTopicsResponse{
	//Topics: CreateTopicResult {
	// Name: "Sounds"
	// ErrorCode: ErrNone
	// ErrorMsg: ""
	//}
	//}

	return resp
}
