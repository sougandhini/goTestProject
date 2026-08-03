package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	// this channel is a key in RSS doc
	// XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"` //each item is a new blog post
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// RSS = structured data in XML format (parsing XML is as same as parsing JSON, XML is just a frappy JSON)
func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body) // This reads everything from resp.Body
	if err != nil {
		return RSSFeed{}, err
	}

	rssfeed := RSSFeed{}

	err = xml.Unmarshal(dat, &rssfeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssfeed, nil

}
