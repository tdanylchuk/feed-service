package service

import (
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
)

type FeedService interface {
	SaveFeed(feed models.Feed) error
	RetrieveFeed(actor string) (*[]models.Feed, error)
	RetrieveFriendsFeed(actor string) (*[]models.Feed, error)
	ProcessAction(actor string, request models.ActionRequest) error
}

func CreateFeedService(
	feedRepository repository.FeedRepository,
	relationRepository repository.RelationRepository) FeedService {
	return &DefaultFeedService{FeedRepository: feedRepository, RelationRepository: relationRepository}
}
