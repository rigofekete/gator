package main

import (
	"fmt"
	"time"
	"context"

	"github.com/rigofekete/gator/internal/database"
	"github.com/google/uuid"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)  
	}

	name := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), name) 
	if err != nil {
		return fmt.Errorf("%s doesn't exist in the database\n", name)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("new username could not be set: %w", err)  
	}

	fmt.Println("User switched successfully")
	return nil 
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)  
	}
	
	name := cmd.Args[0]

	user, err := s.db.CreateUser(
		context.Background(), 
		database.CreateUserParams{
			ID: 	   uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name: 	   name,
		},
	)

	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}


	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("new username could not be set: %w", err)  
	}

	fmt.Printf("User %s has been successfully created\n", name)
	return nil
}


func handlerListUsers(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		} 			
		fmt.Printf("* %s\n", user.Name)
	}

	return nil 
}

