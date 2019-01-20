package models

import (
	"github.com/tdanylchuk/feed-service/entity"
	"time"
)

type FeedResponse struct {
	Actor    string         `json:"actor,omitempty"`
	Verb     string         `json:"verb,omitempty"`
	Object   string         `json:"object,omitempty"`
	Target   string         `json:"target,omitempty"`
	Datetime time.Time      `json:"datetime,omitempty"`
	Related  []FeedResponse `json:"related,omitempty"`
}

func ToResponseFeed(feed entity.FeedEntity) FeedResponse {
	return FeedResponse{
		Actor:    feed.Actor,
		Verb:     feed.Verb,
		Object:   feed.Object,
		Target:   feed.Target,
		Datetime: feed.Datetime,
	}
}
