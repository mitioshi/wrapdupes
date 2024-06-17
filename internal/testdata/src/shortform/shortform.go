package shortform

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidArgument = errors.New("invalid argument")

func SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return errors.New("sql: no rows in result set")
}

type User struct {
	Valid bool
}

func (u User) Validate() error {
	return errors.New("user is in invalid state")
}

func twoWraps() error {
	var user User
	err := SelectContext(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNotFound, err)
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	return nil
}

func wrapWithFormatting() error {
	var user User
	err := SelectContext(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotFound, err.Error())
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidArgument, err.Error())
	}

	return nil
}

func wrapWithUninterpretedString() error {
	var user User
	err := SelectContext(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	}

	return nil
}
