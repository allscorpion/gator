package commands

import (
	"context"
	"fmt"
)

func HandleAgg(s *State, cmd Command) error {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("unable to fetch feed")
	}

	fmt.Println(feed)
	return nil;
}
