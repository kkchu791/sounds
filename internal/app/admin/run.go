package admin

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kkchu791/sounds/internal/domain/service"
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

func Run(flags []string) error {
	servers := os.Getenv("BOOTSTRAP_SERVERS") // "localhost:5001", "localhost:5002"
	fmt.Println(os.Args[5])
	tn := os.Args[4]
	pc, _ := strconv.Atoi(os.Args[6])
	rf, _ := strconv.Atoi(os.Args[8])
	bAddr := strings.Split(servers, ",")[0]
	ctx := context.Background() //TODO: Wrap with Timeout

	c := admin.NewClient(bAddr) // concrete infra client

	topic := service.Topic{
		Name:              tn,
		PartitionCount:    pc,
		ReplicationFactor: rf,
	}

	r, err := service.CreatTopics(ctx, c, bAddr, topic)

	if err != nil {
		fmt.Println(err)
	}

	// r should equal
	// CreateTopicsResponse{
	//Topics: CreateTopicResult {
	// Name: "Sounds"
	// ErrorCode: ErrNone
	// ErrorMsg: ""
	//}
	//}

	fmt.Println(r)
	return nil

	// you have all the env variables now
	// all you have to do is call a service, like create topic
	// service.createtopic

	// then exit this bad boy or print the results of the service

}
