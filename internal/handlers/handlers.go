package handlers

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/scraper"
	"gator/internal/state"
	"net/url"
	"strconv"
	"time"
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

func AddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: gator addFeed <feed name> <feed url>")
	}

	tx, err := s.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	url, err := url.Parse(cmd.Args[1])
	if err != nil {
		return fmt.Errorf("failed to parse URL '%s': %w", cmd.Args[1], err)
	}

	qtx := s.Queries.WithTx(tx)

	feed, err := qtx.CreateFeed(context.Background(),
		database.CreateFeedParams{
			Name:   cmd.Args[0],
			Url:    url.String(),
			UserID: user.ID,
		})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	feedFollow, err := qtx.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			FeedID: feed.ID,
			UserID: user.ID,
		})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("feed '%s' with URL '%s' was added and followed by user '%s'\n", feedFollow.Feed.Name, feedFollow.Feed.Url, feedFollow.User.Name)

	return nil
}

func Agg(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator agg <time between reqs>")
	}

	duration := cmd.Args[0]
	timeBetweenReqs, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("failed to parse time duration '%s': %w", duration, err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scraper.ScrapeFeeds(s); err != nil {
			return fmt.Errorf("error scraping feeds: %v\n", err)
		}
	}
}

func Browse(s *state.State, cmd Command, user database.User) error {
	limit := 2

	if len(cmd.Args) > 0 {
		limitStr := cmd.Args[0]

		lim, err := strconv.Atoi(limitStr)
		if err != nil {
			return fmt.Errorf("invalid limit '%v': %w", lim, err)
		}

		if lim > limit {
			limit = lim
		}
	}

	rows, err := s.Queries.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user '%s': %w", user.Name, err)
	}

	for _, row := range rows {
		line := fmt.Sprintf("---\n")
		line += fmt.Sprintf("* Post:\t%s\n", row.Post.Title)
		line += fmt.Sprintf("* URL:\t%s\n", row.Post.Url)
		line += fmt.Sprintf("* User:\t%s\n", row.User.Name)
		line += fmt.Sprintf("* Date:\t%s\n", row.Post.PublishedAt.Format(time.RFC1123Z))
		line += fmt.Sprintf("---\n")
		fmt.Println(line)
	}

	return nil
}

func Feeds(s *state.State, cmd Command) error {
	rows, err := s.Queries.GetAllFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
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

func Follow(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator follow <url>")
	}

	url := cmd.Args[0]
	feed, err := s.Queries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL '%s': %w", url, err)
	}

	user, err := s.Queries.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	row, err := s.Queries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed '%s': %w", url, err)
	}

	fmt.Printf("feed '%s' is now followed by user '%s'\n", row.Feed.Name, row.User.Name)

	return nil
}

func Following(s *state.State, cmd Command) error {
	user, err := s.Queries.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	rows, err := s.Queries.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feeds for user '%s': %w", user.Name, err)
	}

	if len(rows) == 0 {
		fmt.Println("you are not following any feeds")
		return nil
	}

	fmt.Printf("feeds followed by '%s':\n", user.Name)
	for _, row := range rows {
		fmt.Printf(" - %s\n", row.Feed.Name)
	}

	return nil
}

func Login(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator login <username>")
	}

	username := cmd.Args[0]

	user, err := s.Queries.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get user '%s': %w", username, err)
	}

	s.Cfg.SetUser(user.Name)
	fmt.Printf("logged in as %s\n", user.Name)

	return nil
}

func Register(s *state.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator register <username>")
	}

	username := cmd.Args[0]

	ctx := context.Background()
	user, err := s.Queries.CreateUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to create user '%s': %w", username, err)
	}
	fmt.Printf("user '%s' was created\n", user.Name)

	s.Cfg.SetUser(user.Name)
	fmt.Printf("logged in as '%s'\n", user.Name)

	return nil
}

func Reset(s *state.State, cmd Command) error {
	if err := s.Queries.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}

	fmt.Println("all users have been deleted")

	return nil
}

func Unfollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: gator unfollow <url>")
	}

	url := cmd.Args[0]
	feed, err := s.Queries.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL '%s': %w", url, err)
	}

	if err := s.Queries.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}); err != nil {
		return fmt.Errorf("failed to unfollow feed '%s': %w", url, err)
	}

	fmt.Printf("feed '%s' is no longer followed by user '%s'\n", feed.Name, user.Name)

	return nil
}

func Users(s *state.State, cmd Command) error {
	users, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
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
