package main

import (
	"context"
	"io"
	"fmt"
	"time"
	"net/http"
	"html"
	"encoding/xml"
)


type RSSFeed struct {
	Channel struct {
		Title		string  	`xml:"title"`
		Link 	    	string		`xml:"link"`
		Description	string 		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	} `xml:"channel"`
}


type RSSItem struct {
	Title 		string 	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}



func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Errorf("error making new DO request with context: %w", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")


	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("error sending request: %w", resp)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status code error: %v", resp.StatusCode)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading data from resp body: %w", err)
	}


	var rssFeed RSSFeed 

	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}

