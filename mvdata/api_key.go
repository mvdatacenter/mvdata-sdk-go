package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateAPIKey creates a new API key via POST /auth/api-keys.
func (c *Client) CreateAPIKey(ctx context.Context, create *APIKeyCreate) (*APIKey, error) {
	body, err := encodeBody(create)
	if err != nil {
		return nil, fmt.Errorf("marshaling API key: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/auth/api-keys", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result APIKey
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating API key: %w", err)
	}
	return &result, nil
}

// ListAPIKeys lists all API keys for the current user via GET /auth/api-keys.
func (c *Client) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/auth/api-keys", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result []APIKey
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("listing API keys: %w", err)
	}
	return result, nil
}

// DeleteAPIKey deletes an API key by ID via DELETE /auth/api-keys/:id.
func (c *Client) DeleteAPIKey(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/auth/api-keys/"+id, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting API key %q: %w", id, err)
	}
	return nil
}
