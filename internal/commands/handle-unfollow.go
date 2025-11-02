package commands

import (
	"context"
	"fmt"

	"github.com/allscorpion/gator/internal/database"
)

func HandleUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("feed url is required")
	}

	feedUrl := cmd.Arguments[0];

	feed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl);

	if err != nil {
		return fmt.Errorf("unable getting feed: %v", err);
	}


	err = s.Db.DeleteFeedFollowsForUser(context.Background(), database.DeleteFeedFollowsForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error deleting feed follow %v", err)
	}

	fmt.Println("successfully deleted feed follow");

	return nil;
}