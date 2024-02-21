// Package simple illustrates the trivial case when two fmt.Errorf calls have the same message
package simple

import (
	"errors"
	"fmt"
)

var ErrFoobar = errors.New("foobar")

func foo() error {
	return fmt.Errorf("something went wrong: %w", ErrFoobar)
}

func bar() error {
	return fmt.Errorf("something went wrong: %w", ErrFoobar) // want "duplicate message for a wrapped error: \"something went wrong: %w\""
}
