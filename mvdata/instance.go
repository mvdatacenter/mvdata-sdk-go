package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateInstance creates a new instance via POST /vpcs/{vpcName}/instances.
func (c *Client) CreateInstance(ctx context.Context, vpcName string, instance *Instance) (*Instance, error) {
	reqBody := createInstanceRequest{
		Name:              instance.Name,
		InstanceType:      instance.InstanceType,
		AuthorizedKeyName: instance.AuthorizedKeyName,
	}

	body, err := encodeBody(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling instance: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/vpcs/"+vpcName+"/instances", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Instance
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating instance: %w", err)
	}
	return &result, nil
}

// GetInstance reads an instance by listing all instances in a VPC and
// filtering by name, since the API only exposes GET /vpcs/{vpcName}/instances.
func (c *Client) GetInstance(ctx context.Context, vpcName, name string) (*Instance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/vpcs/"+vpcName+"/instances", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var instances []Instance
	if err := c.do(req, &instances); err != nil {
		return nil, fmt.Errorf("reading instance %q: %w", name, err)
	}

	for i := range instances {
		if instances[i].Name == name {
			return &instances[i], nil
		}
	}
	return nil, &NotFoundError{Resource: "instance " + name}
}

// DeleteInstance removes an instance via DELETE /vpcs/{vpcName}/instances/{name}.
func (c *Client) DeleteInstance(ctx context.Context, vpcName, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/vpcs/"+vpcName+"/instances/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting instance %q: %w", name, err)
	}
	return nil
}
