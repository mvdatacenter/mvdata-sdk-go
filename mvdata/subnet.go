package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateSubnet creates a new subnet via POST /subnets.
func (c *Client) CreateSubnet(ctx context.Context, subnet *Subnet) (*Subnet, error) {
	body, err := encodeBody(subnet)
	if err != nil {
		return nil, fmt.Errorf("marshaling subnet: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subnets", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Subnet
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating subnet: %w", err)
	}
	return &result, nil
}

// GetSubnet reads a subnet by name via GET /subnets/:name.
func (c *Client) GetSubnet(ctx context.Context, name string) (*Subnet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subnets/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Subnet
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading subnet %q: %w", name, err)
	}
	return &result, nil
}

// DeleteSubnet removes a subnet via DELETE /subnets/:name.
func (c *Client) DeleteSubnet(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/subnets/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting subnet %q: %w", name, err)
	}
	return nil
}
