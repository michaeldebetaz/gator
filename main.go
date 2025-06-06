package main

import (
	"gator/internal/config"
	"gator/internal/handlers"
	"gator/internal/state"
	"log"
	"os"
)

func main() {
	s := state.State{Config: config.Read()}

	cmds := handlers.Commands{
		Cmds: make(map[string]func(*state.State, handlers.Command) error),
	}

	cmds.Register("login", handlers.HandlerLogin)

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("Usage gator <command name> <args?>")
	}

	cmd := handlers.Command{Name: args[1], Args: args[2:]}

	if err := cmds.Run(&s, cmd); err != nil {
		log.Fatalf("Error executing command \"%s\": %v", cmd.Name, err)
	}
}
