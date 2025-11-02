package commands

import (
	"context"
	"fmt"

	"github.com/allscorpion/gator/internal/database"
)

func HandleFollowing(s *State, cmd Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error getting feed %v", err)
	}

	fmt.Printf("feeds following:\n")
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.Feedname)
	}

	return nil;
}