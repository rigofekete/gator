package main

import (
	"fmt"
	"time"
	"context"
	"log"
	"database/sql"
	"strings"

	"github.com/rigofekete/gator/internal/database"
	"github.com/google/uuid"
)


func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time Duration>", cmd.Name)  
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration argument: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background()) 
	if err != nil {  
		 log.Println("error fetching next feed", err)
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}


func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched: %w", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url) 
	if err != nil {
		log.Printf("Couldn't collect feed %s: %w", feed.Name, err)
		return
	}


	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: 	t,
				Valid:  true,
			}
		}

		postData := database.CreatePostParams{
			ID: 		uuid.New(),
			CreatedAt: 	time.Now().UTC(),
			UpdatedAt: 	time.Now().UTC(),
			Title:		item.Title,
			Url:		item.Link,
			Description: 	sql.NullString{
				String:	item.Description,
				Valid:	true,
			},
			PublishedAt: 	publishedAt,
			FeedID:		feed.ID,
		}
		_, err = db.CreatePost(context.Background(), postData)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	

	log.Printf("Feed '%s' collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	fmt.Println("=========================================================================")
}


