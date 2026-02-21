package mvdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client wraps the MV Data console REST API.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// New creates a Client with the given base URL and bearer token.
func New(baseURL, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
}

// do executes an HTTP request, sets auth and content-type headers,
// and decodes the JSON response into result (if non-nil).
func (c *Client) do(req *http.Request, result any) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{Resource: req.URL.Path}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

// encodeBody marshals v to JSON for use as a request body.
func encodeBody(v any) (io.Reader, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
