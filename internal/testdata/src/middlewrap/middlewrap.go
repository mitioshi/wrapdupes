// Package middlewrap illustrates the case when a wrapped error is in the middle of the Errorf message
package middlewrap

import (
	"errors"
	"fmt"
)

var ErrFoobar = errors.New("foobar")
var ErrEntryNotFound = errors.New("entry not found")

func searchById() error {
	return fmt.Errorf("%w while searching: %w", ErrEntryNotFound, ErrFoobar)
}

func searchByName() error {
	return fmt.Errorf("%w while searching: %w", ErrEntryNotFound, ErrFoobar) // want "duplicate message for a wrapped error: \"%w while searching: %w\""
}
