package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("login handler expects a username argument")  
	}

	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("new username could not be set. error: %v", err)  
	}

	fmt.Printf("User %s has been successfully set to the config\n", cmd.Args[0])
	return nil
}
