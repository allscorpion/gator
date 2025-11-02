package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/allscorpion/gator/internal/database"
	"github.com/google/uuid"
)

func HandleRegister(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("name is required")
	}

	name := cmd.Arguments[0]
	
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		Name: name,
	});

	if err != nil {
		log.Fatal("user with that name already exists")
	}

	err = s.Cnfg.SetUser(user.Name);

	if err != nil {
		return err;
	}

	fmt.Printf("user was created %v\n", user)
	return nil;
}