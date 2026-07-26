package admin

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kkchu791/sounds/internal/infra/http/admin"
)

func Run(flags []string) error {
	servers := os.Getenv("BOOTSTRAP_SERVERS") // "localhost:5001", "localhost:5002"
	fmt.Println(os.Args[5])
	tn := os.Args[4]
	pc, _ := strconv.Atoi(os.Args[6])
	rf, _ := strconv.Atoi(os.Args[8])
	bAddr := strings.Split(servers, ",")[0]
	ctx := context.Background() //TODO: Wrap with Timeout

	ac := admin.NewClient(bAddr)

	resp, err := ac.GetMetadata(ctx)

	if err != nil {
		fmt.Println(err)
	}

	topic := admin.Topic{
		Name:              tn,
		PartitionCount:    pc,
		ReplicationFactor: rf,
	}

	tl := []admin.Topic{topic}

	cc := admin.NewClient(resp.CAddr)
	resp, err = cc.CreateTopics(ctx, tl)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
	return nil
}
