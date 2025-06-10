package middlewares

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/handlers"
	"gator/internal/state"
)

func LoggedIn(handler func(s *state.State, cmd handlers.Command, user database.User) error) func(s *state.State, cmd handlers.Command) error {
	return func(s *state.State, cmd handlers.Command) error {
		username := s.Cfg.CurrentUserName

		user, err := s.Queries.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("failed to get user '%s': %w", username, err)
		}

		return handler(s, cmd, user)
	}
}
