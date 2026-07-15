package admin

import (
	"fmt"
	"os"
	"strings"

	"github.com/kkchu791/sounds/internal/domain/service"
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
	fmt.Println("hey watsup")
	fmt.Println(os.Args[5])
	tn := os.Args[4]
	pc := os.Args[6]
	rf := os.Args[8]
	bAddr := strings.Split(servers, ",")[0]

	r := service.CreatTopic(tn, pc, rf, bAddr)

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
