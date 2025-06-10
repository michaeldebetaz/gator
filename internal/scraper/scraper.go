package scraper

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"gator/internal/state"
	"strings"
	"time"
)

func ScrapeFeeds(s *state.State) error {
	feed, err := s.Queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	if err := s.Queries.MarkFeedAsFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	for _, item := range rssFeed.Channel.Item {
		if item.Link == "" {
			fmt.Printf("skipping post with title '%s': no link provided\n", item.Title)
			continue
		}

		if item.Title == "" {
			fmt.Printf("skipping post with link '%s': no title provided\n", item.Link)
			continue
		}

		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("failed to parse publication date '%s': %w", item.PubDate, err)
		}

		if _, err := s.Queries.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}); err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				fmt.Printf("skipping post with title '%s': url already exists\n", item.Title)
				continue
			} else {
				return fmt.Errorf("failed to create post: %w", err)
			}
		}
	}

	return nil
}
