package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, u string) (*RSSFeed, error) {
	feed := &RSSFeed{}

	request, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return feed, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return feed, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return feed, fmt.Errorf("request resulted in status code %d %s\v", response.StatusCode, response.Status)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return feed, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := xml.Unmarshal(bytes, feed); err != nil {
		return feed, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	link, err := url.Parse(feed.Channel.Link)
	if err != nil {
		return feed, fmt.Errorf("failed to parse channel link: %w", err)
	}
	feed.Channel.Link = link.String()

	return feed, nil
}
