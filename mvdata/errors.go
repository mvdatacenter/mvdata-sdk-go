package mvdata

import "fmt"

// NotFoundError is returned when the API responds with 404.
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}
