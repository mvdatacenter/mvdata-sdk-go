package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// ListInstanceTypes fetches available instance types and pricing via GET /prices.
func (c *Client) ListInstanceTypes(ctx context.Context) ([]InstanceType, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/prices", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result []InstanceType
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("listing instance types: %w", err)
	}
	return result, nil
}
