package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateInstance creates a new instance via POST /instances.
func (c *Client) CreateInstance(ctx context.Context, instance *Instance) (*Instance, error) {
	reqBody := createInstanceRequest{
		Name:              instance.Name,
		VPCName:           instance.VPCName,
		InstanceType:      instance.InstanceType,
		AuthorizedKeyName: instance.AuthorizedKeyName,
	}

	body, err := encodeBody(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling instance: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/instances", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Instance
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating instance: %w", err)
	}
	return &result, nil
}

// GetInstance reads a single instance via GET /instances/{name}.
func (c *Client) GetInstance(ctx context.Context, name string) (*Instance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/instances/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var instance Instance
	if err := c.do(req, &instance); err != nil {
		return nil, fmt.Errorf("reading instance %q: %w", name, err)
	}
	return &instance, nil
}

// DeleteInstance removes an instance via DELETE /instances/{name}.
func (c *Client) DeleteInstance(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/instances/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting instance %q: %w", name, err)
	}
	return nil
}
