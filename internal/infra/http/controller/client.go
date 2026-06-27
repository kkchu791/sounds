package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type PromoteReq struct {
	PartitionID string
	ISR         []string
}

type PromoteResp struct {
	Status string
}

func (c *Client) doRequest(method, path string, body any) ([]byte, error) {
	var buf bytes.Buffer //creates a Buffer object that grows automatically, implements io.Writer
	json.NewEncoder(&buf).Encode(&body)

	req, err := http.NewRequest(method, c.baseURL+path, &buf)
	if err != nil {
		log.Printf("failed to create request: %s", err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dataByteSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return dataByteSlice, nil
}

func (c *Client) Promote(pID string, isr []string) (*PromoteResp, error) {
	payload := PromoteReq{
		PartitionID: pID,
		ISR:         isr,
	}

	dataByteSlice, err := c.doRequest("POST", "/promote", payload)

	if err != nil {
		return nil, err
	}

	var pr PromoteResp
	err = json.Unmarshal(dataByteSlice, &pr)
	if err != nil {
		return nil, err
	}

	return &pr, err
}
