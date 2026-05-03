package cmd

import (
	"errors"

	"github.com/rigofekete/gator/internal/app"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommands map[string]func(*app.State, Command) error
}

func (c *Commands) Run(s *app.State, cmd Command) error {
	if f, exists := c.RegisteredCommands[cmd.Name]; exists {
		return f(s, cmd)
	}
	return errors.New("command name not available")
}

func (c *Commands) Register(name string, f func(*app.State, Command) error) {
	c.RegisteredCommands[name] = f
}
