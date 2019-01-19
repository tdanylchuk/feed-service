package repository

import (
	"github.com/tdanylchuk/feed-service/models"
)

type FeedRepository interface {
	SaveFeed(feed models.Feed) error
	FindFeeds() (*[]models.Feed, error)
}

func CreateFeedRepository() FeedRepository {
	return &InMemFeedRepository{storedFeeds: []models.Feed{}}
}
