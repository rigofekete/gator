package main

import (
	"errors"
)


type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if f, exists := c.registeredCommands[cmd.Name]; exists == true {
		return f(s, cmd)
	}   
	return errors.New("command name not available")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
