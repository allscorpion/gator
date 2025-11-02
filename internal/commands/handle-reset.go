package commands

import (
	"context"
	"fmt"
	"log"
)

func HandleReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background());

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully deleted all users")
	return nil;
}