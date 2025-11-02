package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/allscorpion/gator/internal/database"
	"github.com/google/uuid"
)

func HandleFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("url is required")
	}

	feedUrl := cmd.Arguments[0];

	currentFeed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl)

	if err != nil {
		return fmt.Errorf("error getting feed %v", err)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		FeedID: currentFeed.ID,
		UserID: user.ID,
	});

	if err != nil {
		return fmt.Errorf("error creating feed follow %v", err)
	}

	fmt.Println("feed follow created")
	fmt.Printf("for username: %v\n", user.Name)
	fmt.Printf("for feed: %v\n", currentFeed.Name)

	return nil;
}