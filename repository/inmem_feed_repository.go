package repository

import (
	"github.com/tdanylchuk/feed-service/models"
)

type InMemFeedRepository struct {
	storedFeeds []models.Feed
}

func (repository *InMemFeedRepository) SaveFeed(feed models.Feed) error {
	repository.storedFeeds = append(repository.storedFeeds, feed)
	return nil
}

func (repository *InMemFeedRepository) FindFeeds() (*[]models.Feed, error) {
	return &repository.storedFeeds, nil
}
