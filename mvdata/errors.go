package mvdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNotFound      = errors.New("mvdata: not found")
	ErrForbidden     = errors.New("mvdata: forbidden")
	ErrLimitExceeded = errors.New("mvdata: limit exceeded")
	ErrConflict      = errors.New("mvdata: conflict")
	ErrInvalid       = errors.New("mvdata: invalid")
	ErrUnavailable   = errors.New("mvdata: unavailable")
)

const (
	CodeNotFound      = "not_found"
	CodeForbidden     = "forbidden"
	CodeLimitExceeded = "limit_exceeded"
	CodeConflict      = "conflict"
	CodeInvalid       = "invalid"
	CodeUnavailable   = "unavailable"
)

var codeSentinels = map[string]error{
	CodeNotFound:      ErrNotFound,
	CodeForbidden:     ErrForbidden,
	CodeLimitExceeded: ErrLimitExceeded,
	CodeConflict:      ErrConflict,
	CodeInvalid:       ErrInvalid,
	CodeUnavailable:   ErrUnavailable,
}

// APIError is a failed request. Status is 0 when no response arrived, and Code is empty for a
// failure the console does not classify, such as a 500.
type APIError struct {
	Status  int
	Code    string
	Message string

	Limit *LimitExceededError
	Err   error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return "executing request: " + e.Err.Error()
	}
	return fmt.Sprintf("API returned %d: %s", e.Status, e.Message)
}

func (e *APIError) Is(target error) bool {
	s, ok := codeSentinels[e.Code]
	return ok && s == target
}

func (e *APIError) Unwrap() []error {
	var errs []error
	if e.Err != nil {
		errs = append(errs, e.Err)
	}
	if e.Limit != nil {
		errs = append(errs, e.Limit)
	}
	return errs
}

type LimitExceededError struct {
	Limit   string
	Current int64
	Max     int64
}

func (e *LimitExceededError) Error() string {
	return fmt.Sprintf("%s limit reached: %d/%d", e.Limit, e.Current, e.Max)
}

// NotFoundError is returned when the API responds with 404.
type NotFoundError struct {
	Resource string
	API      *APIError
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

func (e *NotFoundError) Unwrap() error {
	if e.API == nil {
		return nil
	}
	return e.API
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Limit   string `json:"limit"`
	Current *int64 `json:"current"`
	Max     *int64 `json:"max"`
}

// A body without a code field takes its code from the status, so a console that predates the
// field still classifies.
func newAPIError(status int, body []byte) *APIError {
	e := &APIError{Status: status}
	var b errorBody
	if json.Unmarshal(body, &b) == nil {
		e.Code = b.Code
		e.Message = b.Message
	}
	if e.Message == "" {
		e.Message = strings.TrimSpace(string(body))
	}
	if e.Code == "" {
		e.Code = codeForStatus(status)
	}
	if e.Code == CodeLimitExceeded && b.Limit != "" {
		e.Limit = &LimitExceededError{Limit: b.Limit}
		if b.Current != nil {
			e.Limit.Current = *b.Current
		}
		if b.Max != nil {
			e.Limit.Max = *b.Max
		}
	}
	return e
}

func codeForStatus(status int) string {
	switch status {
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusUnauthorized, http.StatusForbidden:
		return CodeForbidden
	case http.StatusConflict:
		return CodeConflict
	case http.StatusBadRequest:
		return CodeInvalid
	case http.StatusServiceUnavailable:
		return CodeUnavailable
	}
	return ""
}
