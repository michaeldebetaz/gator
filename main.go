package main

import (
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/handlers"
	"gator/internal/middlewares"
	"gator/internal/state"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	s := state.State{Cfg: cfg}

	db, err := sql.Open("postgres", s.Cfg.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	s.Db = db
	s.Queries = database.New(db)

	cmds := handlers.Commands{
		Cmds: make(map[string]func(*state.State, handlers.Command) error),
	}

	cmds.Register("login", handlers.Login)
	cmds.Register("register", handlers.Register)
	cmds.Register("reset", handlers.Reset)
	cmds.Register("users", handlers.Users)
	cmds.Register("agg", handlers.Agg)
	cmds.Register("addfeed", middlewares.LoggedIn(handlers.AddFeed))
	cmds.Register("feeds", handlers.Feeds)
	cmds.Register("follow", handlers.Follow)
	cmds.Register("following", handlers.Following)
	cmds.Register("unfollow", middlewares.LoggedIn(handlers.Unfollow))
	cmds.Register("browse", middlewares.LoggedIn(handlers.Browse))

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("Usage: gator <command name> [args]")
	}

	cmd := handlers.Command{Name: args[1], Args: args[2:]}

	if err := cmds.Run(&s, cmd); err != nil {
		log.Fatalf("Error executing command '%s': %v", cmd.Name, err)
	}
}
