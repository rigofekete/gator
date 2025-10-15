package main

import (
	"errors"

	"github.com/rigofekete/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdsMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.cmdsMap[cmd.name]; exists == true {
		return f(s, cmd)
	}
	return errors.New("command name not available")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdsMap[name] = f
}

