package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kkchu791/sounds/internal/domain/service"
)

// type MetadataReq struct {
// 	Topics []TopicInfo `json:"topics"`
// }

// type TopicInfo struct {
// 	Name string
// }

func (c *Client) Metadata(ctx context.Context) (service.MetadataResponse, error) {
	// payload := MetadataReq{
	// }

	dataByteSlice, err := c.rc.Do(ctx, "GET", "/metadata", nil)

	var mr service.MetadataResponse
	err = json.Unmarshal(dataByteSlice, &mr)
	if err != nil {
		return service.MetadataResponse{}, fmt.Errorf("unmarshal metadata response: %w", err)
	}

	return mr, nil
}
