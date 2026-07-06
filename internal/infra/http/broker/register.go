package broker

import (
	"context"
	"fmt"
)

type RegisterRequest struct {
	BrokerID   string `json:"broker_id"`
	BrokerAddr string `json:"broker_addr"`
}

func (c *Client) Register(ctx context.Context, bID, bAddr string) error {
	payload := RegisterRequest{
		BrokerID:   bID,
		BrokerAddr: bAddr,
	}

	_, err := c.rc.Do(ctx, "POST", "/register", payload)

	if err != nil {
		return fmt.Errorf("do request for /register: %w", err)
	}

	return nil

}
