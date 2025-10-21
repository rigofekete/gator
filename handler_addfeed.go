package main

import (
	"fmt"
	"context"
	"time"

	"github.com/rigofekete/gator/internal/database"
	"github.com/google/uuid"
)


func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <feedName> <feedURL>", cmd.Name)  
	}

	feedName := cmd.Args[0] 
	feedURL := cmd.Args[1] 

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	
	feedParams := database.AddFeedParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: 		feedName,
		Url:  		feedURL,
		UserID: 	user.ID,
	}
	
	feed, err := s.db.AddFeed(context.Background(), feedParams)  
	if err != nil {
		return fmt.Errorf("Error adding feed. %w", err) 
	}

	fmt.Println("New feed: ")
	printFeed(feed)
	return nil 
}

func printFeed(feed database.Feed) {
	fmt.Printf("- ID: 		%s\n", feed.ID)
	fmt.Printf("- Created: 		%v\n", feed.CreatedAt)
	fmt.Printf("- Updated: 		%v\n", feed.UpdatedAt)
	fmt.Printf("- Name: 		%s\n", feed.Name)
	fmt.Printf("- URL: 		%s\n", feed.Url)
	fmt.Printf("- UserID: 		%s\n", feed.UserID)
}
