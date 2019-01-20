package models

import "time"

type FeedResponse struct {
	Actor    string    `json:"actor,omitempty"`
	Verb     string    `json:"verb,omitempty"`
	Object   string    `json:"object,omitempty"`
	Target   string    `json:"target,omitempty"`
	Datetime time.Time `json:"datetime,omitempty"`
	Related  []Feed    `json:"related,omitempty"`
}

func ToResponseFeed(feed Feed) FeedResponse {
	return FeedResponse{
		Actor:    feed.Actor,
		Verb:     feed.Verb,
		Object:   feed.Object,
		Target:   feed.Target,
		Datetime: feed.Datetime,
	}
}
