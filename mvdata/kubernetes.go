package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// CreateKubernetesCluster creates a managed K8s cluster via POST /kubernetes. A cluster with no
// region lands in the deployment default.
func (c *Client) CreateKubernetesCluster(ctx context.Context, cluster *KubernetesCluster) (*KubernetesCluster, error) {
	reqBody := createKubernetesClusterRequest{
		Name:             cluster.Name,
		Version:          cluster.Version,
		NodeInstanceType: cluster.NodeInstanceType,
		NodeCount:        cluster.NodeCount,
	}
	if cluster.Region != nil {
		reqBody.Region = cluster.Region.Code
	}

	body, err := encodeBody(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling kubernetes cluster: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/kubernetes", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result KubernetesCluster
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating kubernetes cluster: %w", err)
	}
	return &result, nil
}

// GetKubernetesCluster reads a cluster via GET /kubernetes/:name.
func (c *Client) GetKubernetesCluster(ctx context.Context, name string) (*KubernetesCluster, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/kubernetes/"+name, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result KubernetesCluster
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading kubernetes cluster %q: %w", name, err)
	}
	return &result, nil
}

// UpdateKubernetesCluster patches mutable fields via PATCH /kubernetes/:name.
func (c *Client) UpdateKubernetesCluster(ctx context.Context, name string, update *KubernetesClusterUpdate) (*KubernetesCluster, error) {
	body, err := encodeBody(update)
	if err != nil {
		return nil, fmt.Errorf("marshaling kubernetes cluster update: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.BaseURL+"/kubernetes/"+name, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result KubernetesCluster
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("updating kubernetes cluster %q: %w", name, err)
	}
	return &result, nil
}

// DeleteKubernetesCluster removes a cluster via DELETE /kubernetes/:name.
func (c *Client) DeleteKubernetesCluster(ctx context.Context, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/kubernetes/"+name, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting kubernetes cluster %q: %w", name, err)
	}
	return nil
}
