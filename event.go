package main

import "time"

type Event struct {
	Feed  EventFeed  `json:"feed"`
	Entry EventEntry `json:"entry"`
}

type EventFeed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type EventEntry struct {
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}
