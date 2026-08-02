package main

import (
	"log"
	"time"

	"github.com/sougandhini/rssagg/internal/database"
)

// this is a logn running job... running as long as my server runs
// so the params are: how many diff go routines you wanna check regulary??


// scraping - using a software to automatically extract data from a website
func startScraping(
	db *database.Queries, concurrency int, timeBetweenRequest time.Duration,
){
	log.Printf("Scraping on %v go-routines for every %s iteration ", concurrency, timeBetweenRequest)
	ticker :=time.NewTicker(timeBetweenRequest) // this creates a sync channel
	for ;; <-ticker.C{ // this means execute the body of the loop everytime a new value comes in
		// this means for every timeBetweenRequest a value will be sent across the channel

	}
	
}


