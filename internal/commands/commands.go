package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/allscorpion/gator/internal/config"
	"github.com/allscorpion/gator/internal/database"
)

type State struct {
	Cnfg *config.Config
	Db * database.Queries; 
}

type Command struct {
	Name string
	Arguments []string
}

type Commands struct {
	Options map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.Options[cmd.Name]
	if !exists {
		return fmt.Errorf("command does not exist")
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Options[name] = f;
}

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("unable to get next feed: %v", err)
	};

	err = s.Db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
		},
		ID: feed.ID,
	});

	if err != nil {
		return fmt.Errorf("unable to mark feed as fetched: %v", err)
	}

	feedDetails, err := FetchFeed(context.Background(), feed.Url);

	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	fmt.Printf("successfully found feed %v\n", feed.Url)
	for _, item := range feedDetails.Channel.Item {
		fmt.Println(item.Title);
	}

	return nil;
}