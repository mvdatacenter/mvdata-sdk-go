package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateKey registers a new SSH key via POST /keys.
func (c *Client) CreateKey(ctx context.Context, key *Key) (*Key, error) {
	body, err := encodeBody(key)
	if err != nil {
		return nil, fmt.Errorf("marshaling key: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/keys", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Key
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating key: %w", err)
	}
	return &result, nil
}

// GetKey reads a key by name via GET /keys?name=X.
func (c *Client) GetKey(ctx context.Context, name string) (*Key, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/keys?name="+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Key
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading key %q: %w", name, err)
	}
	return &result, nil
}

// DeleteKey removes a key via DELETE /keys/:name.
func (c *Client) DeleteKey(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/keys/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting key %q: %w", name, err)
	}
	return nil
}
