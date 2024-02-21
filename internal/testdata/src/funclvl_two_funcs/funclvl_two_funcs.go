package funclvl_two_funcs

import (
	"errors"
	"fmt"
)

var ErrDupe = errors.New("error dupe")

func foo() error {
	return fmt.Errorf("foo: %w", ErrDupe)
}

func bar() error {
	return fmt.Errorf("foo: %w", ErrDupe)
}
