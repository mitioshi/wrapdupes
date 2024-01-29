// Package multifile illustrates the case when multiple files in the same package have the same error message
package multifile

import (
	"errors"
	"fmt"
)

var ErrDupe = errors.New("error dupe")

func foo() error {
	return fmt.Errorf("foo: %w", ErrDupe)
}
