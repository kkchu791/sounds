package broker

import (
	"context"
	"fmt"
)

type HeartbeatReq struct {
	BrokerID string `json:"broker_id"`
}

func (c *Client) Heartbeat(ctx context.Context, bID string) error {
	payload := HeartbeatReq{
		BrokerID: bID,
	}

	_, err := c.rc.Do(ctx, "POST", "/heartbeat", payload)

	if err != nil {
		return fmt.Errorf("do request for /heartbeat: %w", err)
	}

	return nil
}
