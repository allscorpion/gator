package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/allscorpion/gator/internal/database"
	"github.com/google/uuid"
)

func HandleAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("name and url is required")
	}

	name := cmd.Arguments[0];
	url := cmd.Arguments[1];

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		Name: name,
		Url: url,
		UserID: user.ID,
	});

	if err != nil {
		return err;
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		FeedID: feed.ID,
		UserID: user.ID,
	});

	if err != nil {
		return err;
	}

	fmt.Printf("feed successfully added: %v\n", feed);

	return nil;
}