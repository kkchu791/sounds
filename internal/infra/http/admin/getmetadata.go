package admin

import (
	"context"
	"encoding/json"
	"fmt"
)

// type MetadataReq struct {
// 	Topics []TopicInfo `json:"topics"`
// }

// type TopicInfo struct {
// 	Name string
// }

type MetadataResponse struct {
	CAddr string `json:"controller_addr"`
}

func (c *Client) GetMetadata(ctx context.Context) (*MetadataResponse, error) {
	// payload := MetadataReq{
	// }

	dataByteSlice, err := c.rc.Do(ctx, "GET", "/metadata", nil)

	fmt.Println(dataByteSlice, "hello")

	var mr MetadataResponse
	err = json.Unmarshal(dataByteSlice, &mr)
	if err != nil {
		return &MetadataResponse{}, fmt.Errorf("unmarshal metadata response: %w", err)
	}

	return &mr, nil
}
