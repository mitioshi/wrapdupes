package funclvl_samefunc_diff_pkgs

import (
	"errors"
	"fmt"
)

var ErrNotDupe = errors.New("error not dupe")

func funcWithSameName() error {
	return fmt.Errorf("foo: %w", ErrNotDupe)
}
