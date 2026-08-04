package admin

import (
	"context"
	"encoding/json"
	"fmt"
)

// request type structs
type CreateTopicsRequest struct {
	Topics []Topic `json:"topics"`
}

type Topic struct {
	Name              string `json:"topic"`
	PartitionCount    int    `json:"partition_count"`
	ReplicationFactor int    `json:"replication_factor"`
}

// response type structs
type CreateTopicsResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Name      string `json:"name"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

func (c *Client) CreateTopics(ctx context.Context, tl []Topic) (*CreateTopicsResponse, error) {
	payload := CreateTopicsRequest{
		Topics: tl,
	}

	resp, err := c.rc.Do(ctx, "POST", "/createTopic", payload)

	if err != nil {
		return &CreateTopicsResponse{}, fmt.Errorf("do request for /createTopic: %w", err)
	}

	var cr CreateTopicsResponse
	err = json.Unmarshal(resp, &cr)

	if err != nil {
		return &CreateTopicsResponse{}, err
	}

	return &cr, nil
}
