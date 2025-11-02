package commands

import (
	"context"
	"fmt"

	"github.com/allscorpion/gator/internal/database"
)

func getCurrentUser(s *State) (database.User, error) {
	user, err := s.Db.GetUser(context.Background(), s.Cnfg.Current_user_name)

	if err != nil {
		return database.User{}, fmt.Errorf("error getting current user")
	}

	return user, nil;
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	inner := func(st *State, c Command) error {
		currentUser, err := getCurrentUser(st)

		if err != nil {
			return err
		}

		return handler(st, c, currentUser)
	}
	return inner
}