package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.clockify.me/api/v1"

// Client is a Clockify REST API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Clockify client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// do executes an HTTP request against the Clockify API.
func (c *Client) do(method, path string, query url.Values, body any, out any) error {
	fullURL := baseURL + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clockify API error %d: %s", resp.StatusCode, string(respBytes))
	}

	if out != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, out); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) get(path string, query url.Values, out any) error {
	return c.do(http.MethodGet, path, query, nil, out)
}

func (c *Client) post(path string, body any, out any) error {
	return c.do(http.MethodPost, path, nil, body, out)
}

func (c *Client) put(path string, body any, out any) error {
	return c.do(http.MethodPut, path, nil, body, out)
}

func (c *Client) patch(path string, body any, out any) error {
	return c.do(http.MethodPatch, path, nil, body, out)
}

func (c *Client) delete(path string, out any) error {
	return c.do(http.MethodDelete, path, nil, nil, out)
}
