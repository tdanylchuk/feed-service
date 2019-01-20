package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
	"log"
	"time"
)

const (
	FollowRelation = "follow"
)

const (
	FollowVerb   = "follow"
	UnfollowVerb = "unfollow"
)

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
	if feed.Verb == FollowVerb {
		log.Printf("Since feed[%s] is follow request add corresponding relation.", feed)
		return feedService.RelationRepository.AddRelation(feed.Actor, feed.Target, FollowRelation)
	}
	return nil
}

func (feedService *DefaultFeedService) RetrieveFeed(actor string) (*[]models.Feed, error) {
	return feedService.FeedRepository.FindFeedsByActor(actor)
}

func (feedService *DefaultFeedService) RetrieveFriendsFeed(actor string) (*[]models.Feed, error) {
	targets, err := feedService.RelationRepository.GetTargets(actor, FollowRelation)
	if err != nil {
		return nil, err
	}
	return feedService.FeedRepository.FindFeedsByActors(targets)
}

func (feedService *DefaultFeedService) ProcessAction(actor string, request models.ActionRequest) error {
	log.Printf("Processing action request[%#v] from [%s].", request, actor)
	if request.Follow != nil {
		return feedService.ProcessFollowAction(actor, *request.Follow)
	}
	if request.Unfollow != nil {
		return feedService.ProcessUnfollowAction(actor, *request.Unfollow)
	}
	return errors.New("Cannot process unknown action.")
}

func (feedService *DefaultFeedService) ProcessFollowAction(actor string, target string) error {
	log.Printf("Processing follow [%#v] from [%s].", target, actor)
	err := feedService.RelationRepository.AddRelation(actor, target, FollowRelation)
	if err != nil {
		return err
	}

	feed := models.Feed{
		Target: target,
		Actor:  actor,
		Verb:   FollowVerb,
	}
	return feedService.SaveFeed(feed)
}

func (feedService *DefaultFeedService) ProcessUnfollowAction(actor string, target string) error {
	log.Printf("Processing unfollow [%s] from [%s].", target, actor)
	count, err := feedService.RelationRepository.RemoveRelation(actor, target, FollowRelation)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(fmt.Sprintf("Cannot process unfollow request, since [%s] is not following [%s].",
			actor, target))
	}

	feed := models.Feed{
		Target: target,
		Actor:  actor,
		Verb:   UnfollowVerb,
	}
	return feedService.SaveFeed(feed)
}
