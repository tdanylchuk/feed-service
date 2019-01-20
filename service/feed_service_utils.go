package service

import (
	"github.com/tdanylchuk/feed-service/models"
	"reflect"
)

const (
	FollowRelation = "follow"
)

const (
	FollowVerb   = "follow"
	UnfollowVerb = "unfollow"
)

func ConvertToResponseFeeds(feeds *[]models.Feed) *[]models.FeedResponse {
	feedResponses := make([]models.FeedResponse, len(*feeds))
	for i, feed := range *feeds {
		feedResponses[i] = models.ToResponseFeed(feed)
	}
	return &feedResponses
}

func ConvertToResponseFeedsAndCollectUniqueObjects(feeds *[]models.Feed) (*[]models.FeedResponse, *[]string) {
	objectsMap := make(map[string]bool)
	feedResponses := make([]models.FeedResponse, len(*feeds))
	for i, feed := range *feeds {
		feedResponses[i] = models.ToResponseFeed(feed)
		if &feed.Object != nil && len(feed.Object) > 0 {
			objectsMap[feed.Object] = true
		}
	}

	keys := reflect.ValueOf(objectsMap).MapKeys()
	objects := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		objects[i] = keys[i].String()
	}
	return &feedResponses, &objects
}

func EnrichFeedsWithRelated(feedResponses *[]models.FeedResponse, relatedFeeds *[]models.Feed) *[]models.FeedResponse {
	objectFeedsMap := make(map[string][]models.Feed)
	for _, feed := range *relatedFeeds {
		feedList := objectFeedsMap[feed.Object]
		objectFeedsMap[feed.Object] = append(feedList, feed)
	}

	for i := 0; i < len(*feedResponses); i++ {
		feed := &(*feedResponses)[i]
		feed.Related = objectFeedsMap[feed.Object]
	}

	return feedResponses
}
