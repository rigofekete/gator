package main

import (
	"fmt"
	"context"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)  
	}


	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error reseting users database: %w", err) 
	}


	fmt.Println("Database successfully reset")
	return nil 
}
