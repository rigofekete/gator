package main 

import (
	"fmt"
	"context"
	"strconv"

	"github.com/rigofekete/gator/internal/database"
)


func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if limit64, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = limit64
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	} 

	postsData := database.GetPostsForUsersParams{
		UserID:	user.ID,
		Limit:	int32(limit),
	}

	posts, err := s.db.GetPostsForUsers(context.Background(), postsData)
	if err != nil {
		return fmt.Errorf("Error getting posts for user: %w", err)
	}


	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		// Format("Mon Jan 2") stands for week short name, month short name and day of month 
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		// remember Description is an sql.NullString wrapper with String and Valid members
		fmt.Printf("	%v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=========================")
	}


	return nil
}

