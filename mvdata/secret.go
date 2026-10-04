package mvdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const redacted = "[redacted]"

// The string sits behind a pointer so that where fmt cannot call String — a verb such as %d, or a
// Value reached through an unexported field — it prints an address rather than the string.
type Value struct {
	s *string
}

func NewValue(s string) Value {
	return Value{s: &s}
}

func (v Value) Reveal() string {
	if v.s == nil {
		return ""
	}
	return *v.s
}

func (Value) String() string   { return redacted }
func (Value) GoString() string { return redacted }

func (Value) MarshalJSON() ([]byte, error) { return json.Marshal(redacted) }

func (v *Value) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v.s = &s
	return nil
}

func (c *Client) ListSecrets(ctx context.Context, storeName string) ([]SecretMetadata, error) {
	endpoint := c.BaseURL + "/secrets?storeName=" + url.QueryEscape(storeName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result struct {
		Secrets []SecretMetadata `json:"secrets"`
	}
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("listing secrets in store %q: %w", storeName, err)
	}
	return result.Secrets, nil
}

func (c *Client) CreateSecret(ctx context.Context, storeName, name string, value Value) (*SecretMetadata, error) {
	body, err := encodeBody(createSecretRequest{StoreName: storeName, Name: name, Value: value.Reveal()})
	if err != nil {
		return nil, fmt.Errorf("marshaling secret: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/secrets", body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretMetadata
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("creating secret %s/%s: %w", storeName, name, err)
	}
	return &result, nil
}

func (c *Client) GetSecret(ctx context.Context, storeName, name string) (*SecretMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, secretURL(c, storeName, name), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretMetadata
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("reading secret %s/%s: %w", storeName, name, err)
	}
	return &result, nil
}

func (c *Client) UpdateSecretValue(ctx context.Context, storeName, name string, value Value) (*SecretMetadata, error) {
	body, err := encodeBody(updateSecretValueRequest{Value: value.Reveal()})
	if err != nil {
		return nil, fmt.Errorf("marshaling secret value: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, secretURL(c, storeName, name), body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretMetadata
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("updating secret %s/%s: %w", storeName, name, err)
	}
	return &result, nil
}

func (c *Client) DeleteSecret(ctx context.Context, storeName, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, secretURL(c, storeName, name), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if err := c.do(req, nil); err != nil {
		return fmt.Errorf("deleting secret %s/%s: %w", storeName, name, err)
	}
	return nil
}

// The console records each reveal in the account's audit trail, which no other call writes to.
func (c *Client) RevealSecret(ctx context.Context, storeName, name string) (*SecretValue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, secretURL(c, storeName, name)+"/reveal", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	var result SecretValue
	if err := c.do(req, &result); err != nil {
		return nil, fmt.Errorf("revealing secret %s/%s: %w", storeName, name, err)
	}
	return &result, nil
}

func secretURL(c *Client, storeName, name string) string {
	return c.BaseURL + "/secrets/" + url.PathEscape(storeName) + "/" + url.PathEscape(name)
}
