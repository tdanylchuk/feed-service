package service

import (
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
	"log"
	"time"
)

type DefaultFeedService struct {
	FeedRepository repository.FeedRepository
}

func (feedService *DefaultFeedService) SaveFeed(feed models.Feed) error {
	log.Println("Saving new feed.", feed)
	feed.Datetime = time.Now()
	err := feedService.FeedRepository.SaveFeed(feed)
	log.Println("Feed has been saved.", feed)
	return err
}

func (feedService *DefaultFeedService) RetrieveFeed() (*[]models.Feed, error) {
	return feedService.FeedRepository.FindFeeds()
}
