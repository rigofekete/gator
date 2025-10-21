package main

import (
	"fmt"
	"context"
)


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)  
	}

	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml") 
	if err != nil {
		return fmt.Errorf("Error fetching feed. %w", err) 
	}

	fmt.Printf("Feed: %+v", rssFeed)
	return nil 
}
