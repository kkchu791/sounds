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
	servers := os.Getenv("BOOTSTRAP_SERVERS") // "localhost:9001", "localhost:9002"
	fmt.Println(os.Args[5])
	tn := os.Args[4]
	pc, _ := strconv.Atoi(os.Args[6])
	rf, _ := strconv.Atoi(os.Args[8])
	bAddr := strings.Split(servers, ",")[0]
	ctx := context.Background() //TODO: Wrap with Timeout

	ab := admin.NewClient(bAddr)

	resp, err := ab.GetMetadata(ctx)

	if err != nil {
		fmt.Println(err)
	}

	topic := admin.Topic{
		Name:              tn,
		PartitionCount:    pc,
		ReplicationFactor: rf,
	}

	tl := []admin.Topic{topic}

	ac := admin.NewClient(resp.CAddr)
	createResp, err := ac.CreateTopics(ctx, tl)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(createResp, "final thing to look for")
	return nil
}
