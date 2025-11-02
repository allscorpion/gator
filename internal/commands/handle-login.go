package commands

import (
	"context"
	"fmt"
	"log"
)

func HandleLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Arguments[0];

	user, err := s.Db.GetUser(context.Background(), username);

	if err != nil {
		log.Fatal(err)
	}

	err = s.Cnfg.SetUser(user.Name)

	if err != nil {
		return err;
	}

	fmt.Println("username has been set")
	return nil;
}