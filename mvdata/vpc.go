package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateVPC creates a new VPC via POST /vpcs.
func (c *Client) CreateVPC(ctx context.Context, vpc *VPC) (*VPC, error) {
	body, err := encodeBody(vpc)
	if err != nil {
		return nil, fmt.Errorf("marshaling VPC: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/vpcs", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result VPC
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating VPC: %w", err)
	}
	return &result, nil
}

// GetVPC reads a VPC by name via GET /vpcs/:name.
func (c *Client) GetVPC(ctx context.Context, name string) (*VPC, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/vpcs/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result VPC
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading VPC %q: %w", name, err)
	}
	return &result, nil
}

// DeleteVPC removes a VPC via DELETE /vpcs/:name.
func (c *Client) DeleteVPC(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/vpcs/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting VPC %q: %w", name, err)
	}
	return nil
}
