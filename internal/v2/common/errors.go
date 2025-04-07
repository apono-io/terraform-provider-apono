package common

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFoundByName is returned when a resource cannot be found using its name.
	ErrNotFoundByName = errors.New("resource not found by name")
)

func NewNotFoundByNameError(resourceType, name string) error {
	return fmt.Errorf("%s with name '%s' not found: %w", resourceType, name, ErrNotFoundByName)
}
