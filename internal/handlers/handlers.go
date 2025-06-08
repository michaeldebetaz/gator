package handlers

import (
	"context"
	"fmt"
	"gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmds map[string]func(*state.State, Command) error
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	return c.Cmds[cmd.Name](s, cmd)
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.Cmds[name] = f
}

func Login(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("username required")
	}

	username := cmd.Args[0]

	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("failed to get user '%s': %v\n", username, err)
	}

	s.Cfg.SetUser(user.Name)
	fmt.Printf("logged in as %s\n", user.Name)

	return nil
}

func Register(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("name required")
	}

	username := cmd.Args[0]

	ctx := context.Background()
	user, err := s.Db.CreateUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to create user '%s': %v\n", username, err)
	}
	fmt.Printf("user '%s' was created\n", user.Name)

	s.Cfg.SetUser(user.Name)
	fmt.Printf("logged in as '%s'\n", user.Name)

	return nil
}

func Reset(s *state.State, cmd Command) error {
	if err := s.Db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to reset users: %v\n", err)
	}

	fmt.Println("all users have been deleted")

	return nil
}

func Users(s *state.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	for _, user := range users {
		line := fmt.Sprintf("* %s", user.Name)

		if user.Name == s.Cfg.CurrentUserName {
			line += " (current)"
		}

		fmt.Println(line)
	}

	return nil
}
