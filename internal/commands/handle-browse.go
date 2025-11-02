package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/allscorpion/gator/internal/database"
)

func HandleBrowse(s *State, cmd Command, user database.User) error {
	var limit int32;
	limit = 2;

	if len(cmd.Arguments) > 0 {
		if v, err := strconv.Atoi(cmd.Arguments[0]); err == nil {
            limit = int32(v)
        }
	}

	fmt.Println("Getting feeds for user");

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit,
	});

	if err != nil {
		return err;
	}

	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil;
}