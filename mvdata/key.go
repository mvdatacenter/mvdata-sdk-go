package mvdata

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

// GetKey reads a key by name via GET /keys?name=X. The console API has no
// /keys/:name route: it filters the key list by the name query and always
// responds with a JSON array, so the result is decoded as a list and the
// single match is returned. A missing name yields a 404 (NotFoundError) from
// the API; an empty array is treated the same way.
func (c *Client) GetKey(ctx context.Context, name string) (*Key, error) {
	endpoint := c.BaseURL + "/keys?name=" + url.QueryEscape(name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result []Key
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading key %q: %w", name, err)
	}
	if len(result) == 0 {
		return nil, &NotFoundError{Resource: "key " + name}
	}
	return &result[0], nil
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
