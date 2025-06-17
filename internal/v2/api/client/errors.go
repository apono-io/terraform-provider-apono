package client

import (
	"errors"
)

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "404 Not Found"
}

func IsNotFoundError(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}
