package handlers

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
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
	f, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command '%s'", cmd.Name)
	}
	return f(s, cmd)
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

func Agg(s *state.State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	fmt.Println(feed)

	return nil
}

func AddFeed(s *state.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addFeed <feed name> <feed url>")
	}

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	params := database.CreateFeedParams{
		Name:   cmd.Args[0],
		Url:    cmd.Args[1],
		UserID: user.ID,
	}
	feed, err := s.Db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}

	fmt.Println(feed)

	return nil
}

func Feeds(s *state.State, cmd Command) error {
	rows, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %v", err)
	}

	for _, row := range rows {
		line := fmt.Sprintf("---\n")
		line += fmt.Sprintf("* Feed:\t%s\n", row.Feed.Name)
		line += fmt.Sprintf("* URL:\t%s\n", row.Feed.Url)
		line += fmt.Sprintf("* User:\t%s\n", row.User.Name)
		line += fmt.Sprintf("---\n")
		fmt.Println(line)
	}

	return nil
}
