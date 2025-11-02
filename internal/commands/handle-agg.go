package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/allscorpion/gator/internal/database"
)

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

	fmt.Println("----------------------------------------")

	feedDetails, err := FetchFeed(context.Background(), feed.Url);

	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	fmt.Printf("successfully found feed %v\n", feed.Url)
	for _, item := range feedDetails.Channel.Item {
		fmt.Println(item.Title);
	}
	fmt.Println("----------------------------------------")
	fmt.Println("End of feed")


	return nil;
}

func HandleAgg(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("time_between_reqs arg is required");
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0]);

	if err != nil {
		return fmt.Errorf("failed to parse time_between_reqs: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests);

	ticker := time.NewTicker(timeBetweenRequests);

	for ; ; <-ticker.C {
		scrapeFeeds(s);
	}
	
	return nil;
}
