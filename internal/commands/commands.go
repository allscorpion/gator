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

func HandleAgg(s *State, cmd Command) error {
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("unable to fetch feed")
	}

	fmt.Println(feed)
	return nil;
}

func getCurrentUser(s *State) (database.User, error) {
	user, err := s.Db.GetUser(context.Background(), s.Cnfg.Current_user_name)

	if err != nil {
		return database.User{}, fmt.Errorf("error getting current user")
	}

	return user, nil;
}

func HandleAddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("name and url is required")
	}

	name := cmd.Arguments[0];
	url := cmd.Arguments[1];

	currentUser, err := getCurrentUser(s);

	if err != nil {
		return err;
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		Name: name,
		Url: url,
		UserID: currentUser.ID,
	});

	if err != nil {
		return err;
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		FeedID: feed.ID,
		UserID: currentUser.ID,
	});

	if err != nil {
		return err;
	}

	fmt.Printf("feed successfully added: %v\n", feed);

	return nil;
}

func HandlePrintFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
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

func HandleFollowFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("url is required")
	}

	feedUrl := cmd.Arguments[0];
	currentUser, err := getCurrentUser(s)

	if err != nil {
		return err;
	}

	currentFeed, err := s.Db.GetFeedByUrl(context.Background(), feedUrl)

	if err != nil {
		return fmt.Errorf("error getting feed %v", err)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: sql.NullTime{
			Time: time.Now(),
		},
		UpdatedAt: sql.NullTime{
			Time: time.Now(),
		},
		FeedID: currentFeed.ID,
		UserID: currentUser.ID,
	});

	if err != nil {
		return fmt.Errorf("error creating feed follow %v", err)
	}

	fmt.Println("feed follow created")
	fmt.Printf("for username: %v\n", currentUser.Name)
	fmt.Printf("for feed: %v\n", currentFeed.Name)

	return nil;
}

func HandleFollowing(s *State, cmd Command) error {
	currentUser, err := getCurrentUser(s)

	if err != nil {
		return err;
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)

	if err != nil {
		return fmt.Errorf("error getting feed %v", err)
	}

	fmt.Printf("feeds following:\n")
	for _, feedFollow := range feedFollows {
		fmt.Println(feedFollow.Feedname)
	}

	return nil;
}