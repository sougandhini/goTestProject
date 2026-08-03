package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sougandhini/rssagg/internal/database"
)

// this is a logn running job... running as long as my server runs
// so the params are: how many diff go routines you wanna check regulary??

// scraping - using a software to automatically extract data from a website
func startScraping(
	db *database.Queries,
	concurrency int, // how many different go routines on which we wanna do the scraping on
	timeBetweenRequest time.Duration, // time between each request to fetch new RSS feed
) {
	log.Printf("Scraping on %v go-routines for every %s duration ", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest) // this creates a sync channel
	for ; ; <-ticker.C {                         // this means execute the body of the loop everytime a new value comes in
		// this means for every timeBetweenRequest a value will be sent across the channel, go grab concurreny # of feeds at once
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error in fetching feeds: ", err)
			// we continue error, there's no time we want this function to stop
		}

		// fetch feeds individually at the same time, so we use sync mechanism
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}

}
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error in marking feed as fetched: ", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	// Something to note here is that - for each item present in the feed we create a DB entry if its not there previously
	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err: %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post: ", err)
		}

	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
