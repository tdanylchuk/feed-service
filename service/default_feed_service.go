package service

import (
	"github.com/pkg/errors"
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
	"log"
	"time"
)

const followRelation = "follow"

type DefaultFeedService struct {
	FeedRepository     repository.FeedRepository
	RelationRepository repository.RelationRepository
}

func (feedService *DefaultFeedService) SaveFeed(feed models.Feed) error {
	log.Println("Processing new feed request.", feed)
	feed.Datetime = time.Now()
	err := feedService.FeedRepository.SaveFeed(feed)
	log.Println("Feed has been processed.", feed)
	return err
}

func (feedService *DefaultFeedService) ProcessFeed(feed models.Feed) error {
	//feed processor map could added in future map[Verb]Processor
	if feed.Verb == "follow" {
		log.Printf("Since feed[%s] is follow request add corresponding relation.", feed)
		return feedService.RelationRepository.AddRelation(feed.Actor, feed.Target, followRelation)
	}
	return nil
}

func (feedService *DefaultFeedService) RetrieveFeed(actor string) (*[]models.Feed, error) {
	return feedService.FeedRepository.FindFeeds(actor)
}

func (feedService *DefaultFeedService) ProcessAction(actor string, request models.ActionRequest) error {
	log.Printf("Processing action request[%s] from [%s].", request, actor)
	if request.Follow != nil {
		return feedService.ProcessFollowAction(actor, request)
	}
	return errors.New("Cannot process unknown action.")
}

func (feedService *DefaultFeedService) ProcessFollowAction(actor string, request models.ActionRequest) error {
	log.Printf("Processing follow request[%#v] from [%s].", &request, actor)

	//transaction could be applied here to rollback relation if feed processing failed
	err := feedService.RelationRepository.AddRelation(actor, *request.Follow, followRelation)
	if err != nil {
		return err
	}

	feed := models.Feed{
		Target: *request.Follow,
		Actor:  actor,
		Verb:   followRelation,
	}
	return feedService.SaveFeed(feed)
}
