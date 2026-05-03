package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rigofekete/gator/internal/app"
	"github.com/rigofekete/gator/internal/cmd"
	"github.com/rigofekete/gator/internal/database"
)

func HandlerLogin(s *app.State, cmd cmd.Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("%s doesn't exist in the database\n", name)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("new username could not be set: %w", err)
	}

	fmt.Println("User switched successfully")
	return nil
}

func HandlerRegister(s *app.State, cmd cmd.Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.DB.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.Cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("new username could not be set: %w", err)
	}

	fmt.Printf("User %s has been successfully created\n", name)
	return nil
}

func HandlerListUsers(s *app.State, cmd cmd.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
