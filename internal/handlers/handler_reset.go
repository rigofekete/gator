package handlers

import (
	"context"
	"fmt"

	"github.com/rigofekete/gator/internal/app"
	"github.com/rigofekete/gator/internal/cmd"
)

func HandlerReset(s *app.State, cmd cmd.Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error reseting users database: %w", err)
	}

	fmt.Println("Database successfully reset")
	return nil
}
