package mvdata

import (
	"context"
	"fmt"
	"net/http"
)

// DeviceAuthorize initiates the device authorization flow via POST /auth/device/authorize.
func (c *Client) DeviceAuthorize(ctx context.Context) (*DeviceAuth, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/auth/device/authorize", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result DeviceAuth
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("requesting device authorization: %w", err)
	}
	return &result, nil
}

// deviceTokenRequest is the request body for POST /auth/device/token.
type deviceTokenRequest struct {
	DeviceCode string `json:"device_code"`
}

// DeviceToken polls for the device token via POST /auth/device/token.
func (c *Client) DeviceToken(ctx context.Context, deviceCode string) (*DeviceTokenResponse, error) {
	body, err := encodeBody(deviceTokenRequest{DeviceCode: deviceCode})
	if err != nil {
		return nil, fmt.Errorf("marshaling device token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/auth/device/token", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result DeviceTokenResponse
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("polling device token: %w", err)
	}
	return &result, nil
}
