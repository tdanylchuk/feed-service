package controller

import (
	"github.com/tdanylchuk/feed-service/service"
	"net/http"
)

type FeedController interface {
	SaveFeed(w http.ResponseWriter, r *http.Request)
	GetFeeds(w http.ResponseWriter, r *http.Request)
	PerformAction(w http.ResponseWriter, r *http.Request)
}

func CreateController(feedService service.FeedService) FeedController {
	return &DefaultFeedController{FeedService: feedService}
}
