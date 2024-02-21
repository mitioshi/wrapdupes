
# wrapdupes

wrapdupes is a linter to detect duplicate error wraps passed to `fmt.Errorf`

## Motivation

Consider the following example
```go
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
	return fmt.Errorf("something went wrong: %w", ErrFoobar)
}
```
In the logs, you can't tell if foo or bar failed with some error, which defies the purpose of using Errorf in the first place.

Here's another way of adding context to the error you could often see in practice (taken from [mattermost](https://github.com/mattermost/mattermost/blob/v9.4.2/server/channels/app/bot.go#L75))
```go

	// Check for an existing bot user with that username. If one exists, then use that.
	if user, appErr := a.GetUserByUsername(bot.Username); appErr == nil && user != nil {
		if user.IsBot {
			if appErr := a.SetPluginKey(productID, botUserKey, []byte(user.Id)); appErr != nil {
				return "", fmt.Errorf("failed to set plugin key: %w", err)
			}
		} else {
			rctx.Logger().Error("Product attempted to ...", mlog.String("username",
				bot.Username),
				mlog.String("user_id",
					user.Id),
			)
		}
		return user.Id, nil
	}

	createdBot, err := a.CreateBot(rctx, bot)
	if err != nil {
		return "", fmt.Errorf("failed to create bot: %w", err)
	}

	if appErr := a.SetPluginKey(productID, botUserKey, []byte(createdBot.UserId)); appErr != nil {
		return "", fmt.Errorf("failed to set plugin key: %w", err)
	}
```
Here, you can see two distinct code paths that use the the same wrapper for the error: `fmt.Errorf("failed to set plugin key: %w", err)`. If you happened to see this error in logs, you wouldn't be able to tell which code path was taken.

`wrapdupes` should detect such issues in your code.
Both examples are taken from the tests



## License

[MIT](https://choosealicense.com/licenses/mit/)


## Contributing

Contributions are always welcome!

See `contributing.md` for ways to get started.

Please adhere to this project's `code of conduct`.


## Installation

This linter is advised to be used with golangci-lint

You can also install it manually
```bash
go install github.com/mitioshi/wrapdupes/cmd/wrapdupes@latest
```
Then run it by
```bash
wrapdupes ./...
```
    