package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/allscorpion/gator/internal/database"
	"github.com/google/uuid"
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
		parsedTime, err := time.Parse(time.RFC1123, item.PubDate);

		if err != nil {
			return fmt.Errorf("failed to parse the time for %v: %v", item.PubDate, err);
		}

		_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: sql.NullTime{
				Time: time.Now(),
			},
			UpdatedAt: sql.NullTime{
				Time: time.Now(),
			},
			Title: item.Title,
			Url: item.Link,
			Description: item.Description,
			PublishedAt: parsedTime,
			FeedID: feed.ID,
		});

		if err != nil {
			if err.Error() != "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				fmt.Printf("error creating post %v: %v\n", item.Title, err)
			} else {
				continue;
			}
		}

		fmt.Printf("successfully created post for %v\n", item.Title)
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
