package models

import (
	"github.com/tdanylchuk/feed-service/entity"
	"time"
)

type FeedRequest struct {
	Actor  string `json:"actor,omitempty"`
	Verb   string `json:"verb,omitempty"`
	Object string `json:"object,omitempty"`
	Target string `json:"target,omitempty"`
}

func (feedRequest *FeedRequest) ToFeedEntity() entity.FeedEntity {
	return entity.FeedEntity{
		Actor:    feedRequest.Actor,
		Verb:     feedRequest.Verb,
		Object:   feedRequest.Object,
		Target:   feedRequest.Target,
		Datetime: time.Now(),
	}
}
