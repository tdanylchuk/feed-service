package service

import (
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
)

type FeedService interface {
	SaveFeed(feed models.Feed) error
	RetrieveFeed(actor string) (*[]models.Feed, error)
}

func CreateFeedService(feedRepository repository.FeedRepository) FeedService {
	return &DefaultFeedService{FeedRepository: feedRepository}
}
