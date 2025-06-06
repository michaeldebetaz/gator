package handlers

import (
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

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Username required")
	}

	username := cmd.Args[0]
	s.Config.SetUser(username)
	fmt.Printf("Logged in as %s\n", username)

	return nil
}
