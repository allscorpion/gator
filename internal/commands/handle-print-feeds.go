package commands

import (
	"context"
	"fmt"
)

func HandlePrintFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	if len(feeds) == 0 {
		return fmt.Errorf("no feeds found");
	}

	for _, feed := range feeds {
		fmt.Printf("name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)

		if err != nil {
			return fmt.Errorf("unable to get the user details for %v", feed.UserID);
		}

		fmt.Printf("username: %v\n", user.Name)
	}

	return nil;
}