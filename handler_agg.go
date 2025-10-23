package main

import (
	"fmt"
	"time"
	"context"
	"log"

	"github.com/rigofekete/gator/internal/database"
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
		log.Printf("Couldn't mark feed %s fetched: %w", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url) 
	if err != nil {
		log.Printf("Couldn't collect feed %s: %w", feed.Name, err)
		return
	}
	
	fmt.Println("Printing fetched feed items:")
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
	fmt.Println("=================================================")
}


