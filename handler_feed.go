package main

import (
	"fmt"
	"context"
	"time"

	"github.com/rigofekete/gator/internal/database"
	"github.com/google/uuid"
)


func handlerGetFeeds(s* state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err) 
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds registered")
		return nil
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


func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <feedName> <feedURL>", cmd.Name)  
	}

	feedName := cmd.Args[0] 
	feedURL := cmd.Args[1] 

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


	feedFollowParams := database.CreateFeedFollowParams{
		ID:   		uuid.New(), 
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		UserID:    	user.ID,
		FeedID:    	feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams) 
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %w", err)
	}


	fmt.Println("\nNew feed added: ")
	printFeed(feed, user)
	fmt.Println("======================")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("======================")

	return nil 
}


func handlerFollowFeed(s* state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <URL>", cmd.Name)
	}


	url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url) 
	if err != nil {
		return fmt.Errorf("feed not found. %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID: 		uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt:	time.Now().UTC(),
		UserID: 	user.ID,
		FeedID: 	feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}


	fmt.Println("Feed created:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("======================")

	return nil
}


func handlerUnfollowFeed(s* state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <URL>")
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found. %w", err)
	}

	deleteFeedParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), deleteFeedParams)
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %w", err)
	}

	fmt.Printf("%s deleted and unfollowed successfully\n", feed.Name)
	fmt.Println("======================")

	return nil
}

func handlerFollowing(s* state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting current user's feed follows: %w", err)
	}
	
	if len(follows) == 0 {
		fmt.Println("Current user is not following any feed")
		return nil
	}


	fmt.Println("Following: ")
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	fmt.Println("======================")

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


func printFeedFollow(username, feedname string) {
	fmt.Printf("- User Name: %s\n", username)
	fmt.Printf("- Feed Name: %s\n", feedname)
}
