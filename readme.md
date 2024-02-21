
# wrapdupes

wrapdupes is a linter to detect duplicate error wraps passed to `fmt.Errorf` which results in errors in logs being indiscernible from each other.

This linter uses the `go/analysis` API

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
In your logs, you can't tell if foo or bar failed with some error, which defies the purpose of using Errorf in the first place.

Here's another way of adding context to the error you could often see in practice (taken from [mattermost](https://github.com/mattermost/mattermost/blob/v9.4.2/server/channels/app/bot.go#L75))
```

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
## Configuration

`wrapdupes` can be configured to check for duplicate errors either on the package level or a function level
Running `wrapdupes -strictness package ./...` will check if there are any fmt.Errorf calls with exact messages within a **package**.

Running `wrapdupes -strictness function ./...` will check if there are any fmt.Errorf calls with exact messages within a **single function/method**.
If there's two methods with the same wrap message inside a package, this will not be considered an error.

## Examples
Here's what `wrapdupes` outputs for [Mattermost server](https://github.com/mattermost/mattermost/tree/v9.4.2/server)
```
~/src/mattermost/server
wrapdupes ./...
/Users/mitioshi/src/mattermost/server/channels/utils/archive.go:48:17: duplicate message for a wrapped error: "failed to create directory: %w"
/Users/mitioshi/src/mattermost/server/channels/app/imaging/decode.go:85:24: duplicate message for a wrapped error: "imaging: failed to decode image: %w"
/Users/mitioshi/src/mattermost/server/channels/store/sqlstore/shared_channel_store.go:334:15: duplicate message for a wrapped error: "invalid channel: %w"
/Users/mitioshi/src/mattermost/server/channels/store/sqlstore/store.go:1227:18: duplicate message for a wrapped error: "cannot parse MySQL DB version: %w"
/Users/mitioshi/src/mattermost/server/channels/store/sqlstore/store.go:1231:18: duplicate message for a wrapped error: "cannot parse MySQL DB version: %w"
/Users/mitioshi/src/mattermost/server/platform/services/remotecluster/sendprofileImage.go:99:10: duplicate message for a wrapped error: "invalid siteURL while sending file to remote %s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/config.go:157:10: duplicate message for a wrapped error: "invalid config source for %s, %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/service.go:161:16: duplicate message for a wrapped error: "failed to load config from file: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:162:11: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:183:10: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:187:10: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:191:10: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:252:10: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/platform/session.go:259:12: duplicate message for a wrapped error: "%s: %w"
/Users/mitioshi/src/mattermost/server/channels/app/bot.go:100:14: duplicate message for a wrapped error: "failed to set plugin key: %w"
/Users/mitioshi/src/mattermost/server/channels/app/notification_push.go:482:10: duplicate message for a wrapped error: "failed to encode to JSON: %w"
/Users/mitioshi/src/mattermost/server/channels/app/shared_channel.go:43:10: duplicate message for a wrapped error: "cannot find channel: %w"
/Users/mitioshi/src/mattermost/server/channels/app/shared_channel.go:53:11: duplicate message for a wrapped error: "channel is not shared: %w"
/Users/mitioshi/src/mattermost/server/channels/app/shared_channel.go:55:10: duplicate message for a wrapped error: "cannot find channel: %w"
/Users/mitioshi/src/mattermost/server/cmd/mattermost/commands/db.go:245:10: duplicate message for a wrapped error: "failed to initialize filebackend: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/auth.go:226:11: duplicate message for a wrapped error: "could not initiate client: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/auth.go:335:10: duplicate message for a wrapped error: "could not read the password: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/auth.go:344:10: duplicate message for a wrapped error: "could not read the access-token: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/export.go:284:10: duplicate message for a wrapped error: "failed to get export job: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:86:10: duplicate message for a wrapped error: "failed to stat import file: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:102:10: duplicate message for a wrapped error: "failed to create upload session: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:110:10: duplicate message for a wrapped error: "failed to upload data: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:123:10: duplicate message for a wrapped error: "failed to create import process job: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:309:11: duplicate message for a wrapped error: "cannot encode user line: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:317:11: duplicate message for a wrapped error: "cannot encode user line: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:353:12: duplicate message for a wrapped error: "cannot encode post line: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:369:11: duplicate message for a wrapped error: "cannot encode channel line: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/sampledata.go:387:12: duplicate message for a wrapped error: "cannot encode post line: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/utils_unix.go:70:10: duplicate message for a wrapped error: "could not initiate prompt: %w"
/Users/mitioshi/src/mattermost/server/cmd/mmctl/commands/utils_unix.go:74:10: duplicate message for a wrapped error: "error running prompt: %w""
```

