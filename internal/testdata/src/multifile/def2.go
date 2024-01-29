package multifile

import "fmt"

func bar() error {
	return fmt.Errorf("foo: %w", ErrDupe) // want "duplicate message for a wrapped error: \"foo: %w\""
}
