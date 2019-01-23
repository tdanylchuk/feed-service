package service

import (
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
	"github.com/tdanylchuk/feed-service/sender"
)

type FeedService interface {
	SaveFeed(feed models.FeedRequest) error
	ProcessFeed(feed models.FeedRequest) error
	RetrieveFeed(actor string, includeRelated bool) (*[]models.FeedResponse, error)
	RetrieveFriendsFeed(actor string) (*[]models.FeedResponse, error)
	ProcessAction(actor string, request models.ActionRequest) error
}

func CreateFeedService(
	feedRepository repository.FeedRepository,
	relationRepository repository.RelationRepository,
	sender sender.Sender) FeedService {
	return &DefaultFeedService{FeedRepository: feedRepository, RelationRepository: relationRepository, Sender: sender}
}
