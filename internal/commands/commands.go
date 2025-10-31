package commands

import (
	"fmt"

	"github.com/allscorpion/gator/internal/config"
)

type State struct {
	Cnfg *config.Config
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

	username := cmd.Arguments[0]
	err := s.Cnfg.SetUser(username)

	if err != nil {
		return err;
	}

	fmt.Println("username has been set")
	return nil;
}