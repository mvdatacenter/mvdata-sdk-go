package mvdata

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// ListRegions fetches the regions the platform serves. The set is seeded platform data that grows
// as regions are brought up, so a caller reads it rather than carrying its own copy of the codes.
func (c *Client) ListRegions(ctx context.Context) ([]Region, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/regions", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result []Region
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("listing regions: %w", err)
	}
	return result, nil
}

// GetRegion reads a region by code via GET /regions?code=:code. The API answers the filtered query
// with the region itself rather than a list of one.
func (c *Client) GetRegion(ctx context.Context, code string) (*Region, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/regions?code="+url.QueryEscape(code), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result Region
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading region %q: %w", code, err)
	}
	return &result, nil
}
