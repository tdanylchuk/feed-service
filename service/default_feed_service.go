package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/repository"
	"github.com/tdanylchuk/feed-service/sender"
	"log"
)

type DefaultFeedService struct {
	FeedRepository     repository.FeedRepository
	RelationRepository repository.RelationRepository
	Sender             sender.Sender
}

func (feedService *DefaultFeedService) ProcessFeed(feedRequest models.FeedRequest) error {
	log.Println("Processing new feed request.", feedRequest)
	err := feedService.processFeed(feedRequest)
	if err != nil {
		return err
	}
	err = feedService.Sender.Send(&feedRequest)
	log.Println("Feed has been processed.", feedRequest)
	return err
}

func (feedService *DefaultFeedService) SaveFeed(feedRequest models.FeedRequest) error {
	log.Println("Saving new feed ...", feedRequest)
	feed := feedRequest.ToFeedEntity()
	err := feedService.FeedRepository.SaveFeed(feed)
	log.Println("Feed has been saved.", feedRequest)
	return err
}

func (feedService *DefaultFeedService) processFeed(feed models.FeedRequest) error {
	//feed processor map could added in future - map[Verb]Processor
	if feed.Verb == FollowVerb {
		log.Printf("Since feed[%s] is follow request add corresponding relation.", feed)
		return feedService.RelationRepository.AddRelation(feed.Actor, feed.Target, FollowRelation)
	}
	return nil
}

func (feedService *DefaultFeedService) RetrieveFeed(actor string, includeRelated bool) (*[]models.FeedResponse, error) {
	feeds, err := feedService.FeedRepository.FindFeedsByActor(actor)
	if err != nil {
		return nil, err
	}

	if !includeRelated {
		return ConvertToResponseFeeds(feeds), nil
	}

	enrichedFeeds, uniqueObjects := ConvertToResponseFeedsAndCollectUniqueObjects(feeds)
	if len(*uniqueObjects) == 0 {
		return enrichedFeeds, nil
	}

	targets, err := feedService.RelationRepository.GetTargets(actor, FollowRelation)
	if err != nil {
		return nil, err
	}

	if len(*targets) == 0 {
		return enrichedFeeds, nil
	}

	relatedFeeds, err := feedService.FeedRepository.FindFeedsByActorsAndObjects(targets, uniqueObjects)
	if err != nil {
		return nil, err
	}
	return EnrichFeedsWithRelated(enrichedFeeds, ConvertToResponseFeeds(relatedFeeds)), nil
}

func (feedService *DefaultFeedService) RetrieveFriendsFeed(actor string) (*[]models.FeedResponse, error) {
	targets, err := feedService.RelationRepository.GetTargets(actor, FollowRelation)
	if err != nil {
		return nil, err
	}

	feeds, err := feedService.FeedRepository.FindFeedsByActors(targets)
	if err != nil {
		return nil, err
	}
	return ConvertToResponseFeeds(feeds), nil
}

func (feedService *DefaultFeedService) ProcessAction(actor string, request models.ActionRequest) error {
	log.Printf("Processing action request[%#v] from [%s].", request, actor)
	if request.Follow != nil {
		return feedService.processFollowAction(actor, *request.Follow)
	}
	if request.Unfollow != nil {
		return feedService.processUnfollowAction(actor, *request.Unfollow)
	}
	return errors.New("Cannot process unknown action.")
}

func (feedService *DefaultFeedService) processFollowAction(actor string, target string) error {
	log.Printf("Processing follow [%#v] from [%s].", target, actor)
	err := feedService.RelationRepository.AddRelation(actor, target, FollowRelation)
	if err != nil {
		return err
	}

	feed := models.FeedRequest{
		Target: target,
		Actor:  actor,
		Verb:   FollowVerb,
	}
	return feedService.Sender.Send(&feed)
}

func (feedService *DefaultFeedService) processUnfollowAction(actor string, target string) error {
	log.Printf("Processing unfollow [%s] from [%s].", target, actor)
	count, err := feedService.RelationRepository.RemoveRelation(actor, target, FollowRelation)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(fmt.Sprintf("Cannot process unfollow request, since [%s] is not following [%s].",
			actor, target))
	}

	feed := models.FeedRequest{
		Target: target,
		Actor:  actor,
		Verb:   UnfollowVerb,
	}
	return feedService.Sender.Send(&feed)
}
