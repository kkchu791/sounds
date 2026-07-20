package service

import (
	"context"
	"fmt"

	"github.com/kkchu791/sounds/internal/infra/http/admin"
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

type TopicList struct {
	Topics []Topic
}

type Topic struct {
	Name              string
	PartitionCount    int
	ReplicationFactor int
}

func CreatTopics(ctx context.Context, bAddr string, topic ...Topic) (CreateTopicsResponse, error) {
	fmt.Println("in the service, Hey")

	//creates a client for broker
	aClient := admin.NewClient(bAddr)
	resp, err := aClient.Metadata(ctx)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		return CreateTopicsResponse{}, err
	}

	//resp should equal MetaDataRes{controller_addr: "localhost:5000"}

	// cAddr := resp.controller_addr
	// cClient := controller.NewClient(cAddr)

	// passed in a list of topics

	tl := TopicList{
		[]Topic{topic[0]},
	}

	// resp, err = cClient.CreateTopics(ctx, tl)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return CreateTopicsResponse{}, err
	// }
	// resp should equal
	// CreateTopicsResponse{
	//Topics: CreateTopicResult {
	// Name: "Sounds"
	// ErrorCode: ErrNone
	// ErrorMsg: ""
	//}
	//}

	return CreateTopicsResponse{}, err
}
