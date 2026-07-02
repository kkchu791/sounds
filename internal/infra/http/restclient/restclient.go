package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RestClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRestClient(baseURL string) *RestClient {
	return &RestClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (rc *RestClient) Do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var buf bytes.Buffer //creates a Buffer object that grows automatically, implements io.Writer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)

		if err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, rc.baseURL+path, &buf)

	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	defer resp.Body.Close()

	dataByteSlice, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, dataByteSlice)
	}

	return dataByteSlice, nil
}
