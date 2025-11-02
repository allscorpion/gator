package commands

import (
	"context"
	"fmt"
	"log"
)

func HandleGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background());

	if err != nil {
		log.Fatal(err)
	}

	currentUserName := s.Cnfg.Current_user_name

	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil;
}