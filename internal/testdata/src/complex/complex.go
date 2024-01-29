package complex

import (
	"fmt"
)

type Bot struct {
	ID string
}

type User struct {
	ID string
}

func (u User) IsBot() bool {
	return false
}

func setPluginKey(userID string) error {
	return nil
}

func createBot(bot Bot) (Bot, error) {
	return Bot{}, nil
}

func getUserByUsername(name string) (User, error) {
	return User{}, nil
}

func ensureBot(bot Bot) error {
	// Check for an existing bot user with that username. If one exists, then use that.
	if user, err := getUserByUsername("matz"); err == nil {
		if user.IsBot() {
			if err := setPluginKey(user.ID); err != nil {
				return fmt.Errorf("failed to set plugin key: %w", err)
			}
		}
	}
	createdBot, err := createBot(bot)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	if err := setPluginKey(createdBot.ID); err != nil {
		return fmt.Errorf("failed to set plugin key: %w", err) // want "duplicate message for a wrapped error: \"failed to set plugin key: %w\""
	}

	return nil
}
