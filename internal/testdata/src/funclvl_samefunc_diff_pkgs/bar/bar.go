package bar

import (
	"errors"
	"fmt"
)

var ErrDupe = errors.New("error dupe")

func funcWithSameName() error {
	return fmt.Errorf("foo: %w", ErrDupe)
}
