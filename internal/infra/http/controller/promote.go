package controller

import (
	"context"
	"encoding/json"
	"fmt"
)

type PromoteReq struct {
	PartitionID string   `json:"partition_id"`
	ISR         []string `json:"isr"`
}

type PromoteResp struct {
	Status string `json:"status"`
}

func (c *Client) Promote(ctx context.Context, pID string, isr []string) (*PromoteResp, error) {
	payload := PromoteReq{
		PartitionID: pID,
		ISR:         isr,
	}

	dataByteSlice, err := c.rc.Do(ctx, "POST", "/promote", payload)

	if err != nil {
		return nil, fmt.Errorf("do request for /promote: %w", err)
	}

	var pr PromoteResp
	err = json.Unmarshal(dataByteSlice, &pr)
	if err != nil {
		return nil, fmt.Errorf("unmarshal promote response: %w", err)
	}

	return &pr, nil
}
