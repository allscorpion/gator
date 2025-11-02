package commands

import (
	"fmt"

	"github.com/allscorpion/gator/internal/config"
	"github.com/allscorpion/gator/internal/database"
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

