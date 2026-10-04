package mvdata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSecretRoutes(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		call       func(*Client) error
		wantMethod string
		wantURI    string
		wantBody   string
		respond    string
	}{
		{"ListSecretStores", func(c *Client) error {
			got, err := c.ListSecretStores(ctx)
			if err == nil && (len(got) != 1 || got[0].Name != "app" || got[0].SecretCount != 2) {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodGet, "/secret-stores", "", `{"secretStores":[{"name":"app","secretCount":2}]}`},
		{"CreateSecretStore", func(c *Client) error {
			got, err := c.CreateSecretStore(ctx, &SecretStore{Name: "app", Description: "the app"})
			if err == nil && got.Name != "app" {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodPost, "/secret-stores", `{"name":"app","description":"the app"}`, `{"name":"app"}`},
		{"GetSecretStore", func(c *Client) error {
			got, err := c.GetSecretStore(ctx, "app")
			if err == nil && got.Description != "the app" {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodGet, "/secret-stores/app", "", `{"name":"app","description":"the app"}`},
		{"UpdateSecretStore", func(c *Client) error {
			_, err := c.UpdateSecretStore(ctx, "app", &SecretStoreUpdate{})
			return err
		}, http.MethodPatch, "/secret-stores/app", `{"description":""}`, `{"name":"app"}`},
		{"DeleteSecretStore", func(c *Client) error {
			return c.DeleteSecretStore(ctx, "app")
		}, http.MethodDelete, "/secret-stores/app", "", ""},
		{"ListSecrets", func(c *Client) error {
			got, err := c.ListSecrets(ctx, "app")
			if err == nil && (len(got) != 1 || got[0].Version != 3 || got[0].SyncState != "in-sync") {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodGet, "/secrets?storeName=app", "", `{"secrets":[{"storeName":"app","name":"db","version":3,"syncState":"in-sync"}]}`},
		{"CreateSecret", func(c *Client) error {
			got, err := c.CreateSecret(ctx, "app", "db", NewValue("hunter2"))
			if err == nil && got.Version != 1 {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodPost, "/secrets", `{"storeName":"app","name":"db","value":"hunter2"}`, `{"storeName":"app","name":"db","version":1}`},
		{"GetSecret", func(c *Client) error {
			got, err := c.GetSecret(ctx, "app", "db")
			if err == nil && got.Version != 3 {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodGet, "/secrets/app/db", "", `{"storeName":"app","name":"db","version":3}`},
		{"UpdateSecretValue", func(c *Client) error {
			got, err := c.UpdateSecretValue(ctx, "app", "db", NewValue("hunter3"))
			if err == nil && got.Version != 4 {
				err = fmt.Errorf("decoded %+v", got)
			}
			return err
		}, http.MethodPut, "/secrets/app/db", `{"value":"hunter3"}`, `{"storeName":"app","name":"db","version":4}`},
		{"DeleteSecret", func(c *Client) error {
			return c.DeleteSecret(ctx, "app", "db")
		}, http.MethodDelete, "/secrets/app/db", "", ""},
		{"RevealSecret", func(c *Client) error {
			got, err := c.RevealSecret(ctx, "app", "db")
			if err == nil && (got.Value.Reveal() != "hunter2" || got.Version != 4) {
				err = fmt.Errorf("decoded version %d", got.Version)
			}
			return err
		}, http.MethodPost, "/secrets/app/db/reveal", "", `{"storeName":"app","name":"db","version":4,"value":"hunter2"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod, gotURI, gotBody string
			client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotURI = r.URL.RequestURI()
				b, _ := io.ReadAll(r.Body)
				gotBody = string(b)
				if tt.respond == "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.Write([]byte(tt.respond))
			}))

			if err := tt.call(client); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotMethod != tt.wantMethod {
				t.Errorf("expected %s, got %s", tt.wantMethod, gotMethod)
			}
			if gotURI != tt.wantURI {
				t.Errorf("expected %s, got %s", tt.wantURI, gotURI)
			}
			if gotBody != tt.wantBody {
				t.Errorf("expected body %s, got %s", tt.wantBody, gotBody)
			}
		})
	}
}

// A name the console would refuse still stays in its own path segment.
func TestSecretNamesAreEscapedIntoOnePathSegment(t *testing.T) {
	var gotPath string
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.Write([]byte(`{}`))
	}))

	if _, err := client.GetSecret(context.Background(), "app", "../x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotPath != "/secrets/app/..%2Fx" {
		t.Errorf("expected /secrets/app/..%%2Fx, got %s", gotPath)
	}
}

func TestValueNeverFormatsItsString(t *testing.T) {
	v := NewValue("hunter2")
	sv := SecretValue{StoreName: "app", Name: "db", Version: 1, Value: v}
	type wrapper struct{ value Value }

	for _, format := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%X", "%10s", "%d"} {
		for _, arg := range []any{v, &v, sv, &sv, []Value{v}, map[string]Value{"k": v}, wrapper{v}} {
			if got := fmt.Sprintf(format, arg); strings.Contains(got, "hunter2") ||
				strings.Contains(got, fmt.Sprintf("%x", "hunter2")) || strings.Contains(got, fmt.Sprintf("%X", "hunter2")) {
				t.Errorf("Sprintf(%q, %T) = %q, which shows the value", format, arg, got)
			}
		}
	}
	if got := fmt.Sprint(v); got != "[redacted]" {
		t.Errorf("expected [redacted], got %q", got)
	}
	if got := v.Reveal(); got != "hunter2" {
		t.Errorf("expected Reveal to return hunter2, got %q", got)
	}
	if got := (Value{}).Reveal(); got != "" {
		t.Errorf("expected the zero Value to reveal an empty string, got %q", got)
	}
}

func TestValueMarshalsRedactedAndUnmarshalsTheString(t *testing.T) {
	b, err := json.Marshal(SecretValue{Name: "db", Value: NewValue("hunter2")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(b), "hunter2") || !strings.Contains(string(b), `"value":"[redacted]"`) {
		t.Errorf("expected a redacted value, got %s", b)
	}

	var sv SecretValue
	if err := json.Unmarshal([]byte(`{"name":"db","value":"hunter2"}`), &sv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Value.Reveal() != "hunter2" {
		t.Errorf("expected hunter2, got %q", sv.Value.Reveal())
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		wantCode string
		want     error
	}{
		{"coded not_found", 404, `{"statusCode":404,"error":"Not Found","code":"not_found","message":"no such secret"}`, CodeNotFound, ErrNotFound},
		{"coded forbidden", 403, `{"statusCode":403,"error":"Forbidden","code":"forbidden","message":"scope"}`, CodeForbidden, ErrForbidden},
		{"coded limit_exceeded", 403, `{"statusCode":403,"error":"Forbidden","code":"limit_exceeded","limit":"maxSecrets","current":50,"max":50,"message":"Account limit reached: 50/50 secrets."}`, CodeLimitExceeded, ErrLimitExceeded},
		{"coded conflict", 409, `{"statusCode":409,"error":"Conflict","code":"conflict","message":"taken"}`, CodeConflict, ErrConflict},
		{"coded invalid", 400, `{"statusCode":400,"error":"Bad Request","code":"invalid","message":"bad name"}`, CodeInvalid, ErrInvalid},
		{"coded unavailable", 503, `{"statusCode":503,"error":"Service Unavailable","code":"unavailable","message":"pyx"}`, CodeUnavailable, ErrUnavailable},
		{"uncoded 401", 401, `{"statusCode":401,"error":"Unauthorized","message":"Missing authentication"}`, CodeForbidden, ErrForbidden},
		{"uncoded 403", 403, `{"statusCode":403,"error":"Forbidden","message":"scope"}`, CodeForbidden, ErrForbidden},
		{"uncoded 409", 409, `{"statusCode":409,"error":"Conflict","message":"taken"}`, CodeConflict, ErrConflict},
		{"uncoded 400", 400, `{"statusCode":400,"error":"Bad Request","message":"bad"}`, CodeInvalid, ErrInvalid},
		{"uncoded 503", 503, `upstream down`, CodeUnavailable, ErrUnavailable},
		{"uncoded 500", 500, `{"error":"Internal Server Error"}`, "", nil},
	}
	sentinels := []error{ErrNotFound, ErrForbidden, ErrLimitExceeded, ErrConflict, ErrInvalid, ErrUnavailable}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))

			_, err := client.CreateSecret(context.Background(), "app", "db", NewValue("hunter2"))
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected an APIError, got %T: %v", err, err)
			}
			if apiErr.Status != tt.status || apiErr.Code != tt.wantCode {
				t.Errorf("expected status %d code %q, got %d %q", tt.status, tt.wantCode, apiErr.Status, apiErr.Code)
			}
			for _, s := range sentinels {
				if got := errors.Is(err, s); got != (s == tt.want) {
					t.Errorf("errors.Is(err, %v) = %v", s, got)
				}
			}
			if strings.Contains(err.Error(), "hunter2") {
				t.Errorf("the error shows the value: %v", err)
			}
		})
	}
}

func TestErrorMessageIsTheConsoleMessage(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"statusCode":409,"error":"Conflict","code":"conflict","message":"secret store app still holds secrets"}`))
	}))

	err := client.DeleteSecretStore(context.Background(), "app")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected an APIError, got %T: %v", err, err)
	}
	if apiErr.Message != "secret store app still holds secrets" {
		t.Errorf("expected the console's message, got %q", apiErr.Message)
	}
}

func TestLimitExceededCarriesTheBound(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"statusCode":403,"error":"Forbidden","code":"limit_exceeded","limit":"maxSecretKeys","current":200,"max":200,"message":"Account limit reached"}`))
	}))

	_, err := client.UpdateSecretValue(context.Background(), "app", "db", NewValue("x"))
	var limit *LimitExceededError
	if !errors.As(err, &limit) {
		t.Fatalf("expected a LimitExceededError, got %T: %v", err, err)
	}
	if limit.Limit != "maxSecretKeys" || limit.Current != 200 || limit.Max != 200 {
		t.Errorf("expected maxSecretKeys 200/200, got %+v", limit)
	}
}

// Callers that check for *NotFoundError, as the provider does, keep working.
func TestNotFoundKeepsItsTypeAndCarriesTheCode(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"statusCode":404,"error":"Not Found","code":"not_found","message":"no such secret"}`))
	}))

	_, err := client.GetSecret(context.Background(), "app", "db")
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Fatalf("expected a NotFoundError, got %T: %v", err, err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Message != "no such secret" {
		t.Errorf("expected the NotFoundError to carry the console's message, got %v", apiErr)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected errors.Is(err, ErrNotFound)")
	}
}

func TestTransportFailureIsUnavailable(t *testing.T) {
	client := New("http://127.0.0.1:1", "test-token")

	_, err := client.GetSecret(context.Background(), "app", "db")
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %T: %v", err, err)
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != 0 || apiErr.Err == nil {
		t.Errorf("expected an APIError with no status and the transport error, got %+v", apiErr)
	}
}

func TestCanceledRequestStillMatchesContextCanceled(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetSecret(ctx, "app", "db")
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected errors.Is(err, context.Canceled), got %v", err)
	}
}
