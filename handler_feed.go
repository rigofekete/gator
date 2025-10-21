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
	fmt.Println()
	printFeed(feed, user)
	return nil 
}

func handlerGetFeeds(s* state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err) 
	}
	
	fmt.Println("Current feeds: ")
fmt.Println()
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error getting user name from users table: %w", err)
		}
		printFeed(feed, user)
		fmt.Println()
	}
	return nil
}

func handlerFollowFeed(s* state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <URL>", cmd.Name)
	}

	url := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), url) 
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		UserID: 	user.ID,
		FeedID: 	feed.ID,
	}

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("DEBUG LOG: \n %+v", feedFollowRow)

	return nil
}


func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("- ID: 		%s\n", feed.ID)
	fmt.Printf("- Created: 		%v\n", feed.CreatedAt)
	fmt.Printf("- Updated: 		%v\n", feed.UpdatedAt)
	fmt.Printf("- Name: 		%s\n", feed.Name)
	fmt.Printf("- URL: 		%s\n", feed.Url)
	fmt.Printf("- UserID: 		%s\n", feed.UserID)
	fmt.Printf("- User Name: 	%s\n", user.Name)
}
