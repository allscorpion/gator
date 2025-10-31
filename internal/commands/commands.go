package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/allscorpion/gator/internal/config"
	"github.com/allscorpion/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Cnfg *config.Config
	Db * database.Queries; 
}

type Command struct {
	Name string
	Arguments []string
}

type Commands struct {
	Options map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.Options[cmd.Name]
	if !exists {
		return fmt.Errorf("command does not exist")
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Options[name] = f;
}

func HandlerLogin(s *State, cmd Command) error {
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

func HandlerRegister(s *State, cmd Command) error {
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

func HandleReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background());

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully deleted all users")
	return nil;
}

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