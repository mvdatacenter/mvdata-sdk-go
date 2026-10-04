package mvdata

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) ListSecretStores(ctx context.Context) ([]SecretStore, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/secret-stores", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result struct {
		SecretStores []SecretStore `json:"secretStores"`
	}
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("listing secret stores: %w", err)
	}
	return result.SecretStores, nil
}

func (c *Client) CreateSecretStore(ctx context.Context, store *SecretStore) (*SecretStore, error) {
	body, err := encodeBody(createSecretStoreRequest{Name: store.Name, Description: store.Description})
	if err != nil {
		return nil, fmt.Errorf("marshaling secret store: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/secret-stores", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretStore
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating secret store %q: %w", store.Name, err)
	}
	return &result, nil
}

func (c *Client) GetSecretStore(ctx context.Context, name string) (*SecretStore, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, secretStoreURL(c, name), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretStore
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading secret store %q: %w", name, err)
	}
	return &result, nil
}

func (c *Client) UpdateSecretStore(ctx context.Context, name string, update *SecretStoreUpdate) (*SecretStore, error) {
	body, err := encodeBody(update)
	if err != nil {
		return nil, fmt.Errorf("marshaling secret store update: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, secretStoreURL(c, name), body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretStore
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("updating secret store %q: %w", name, err)
	}
	return &result, nil
}

func (c *Client) DeleteSecretStore(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, secretStoreURL(c, name), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting secret store %q: %w", name, err)
	}
	return nil
}

func secretStoreURL(c *Client, name string) string {
	return c.BaseURL + "/secret-stores/" + url.PathEscape(name)
}
