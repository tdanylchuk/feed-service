package controller

import (
	"encoding/json"
	"fmt"
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/service"
	"log"
	"net/http"
)

type DefaultFeedController struct {
	FeedService service.FeedService
}

func (controller *DefaultFeedController) SaveFeed(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller. Saving new feed...")
	actor := GetActor(r)
	var feed models.Feed
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&feed); err != nil {
		log.Println("Something went wrong during saving feed...", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if feed.Actor != actor {
		str := fmt.Sprintf("Actor[%s] is not eligible to post others[%s] feed items.", actor, feed.Actor)
		log.Printf(str)
		respondWithError(w, http.StatusForbidden, str)
		return
	}

	if err := controller.FeedService.SaveFeed(feed); err != nil {
		log.Printf("Something went wrong during feed saving. Feed - [%s]. Error - [%s]", feed, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (controller *DefaultFeedController) GetFeeds(w http.ResponseWriter, r *http.Request) {
	actor := GetActor(r)
	log.Printf("Controller. Retrieving user feed for [%s]...", actor)
	feeds, err := controller.FeedService.RetrieveFeed(actor)
	if err != nil {
		log.Printf("Something went wrong during retriving feeds for [%s]. Error - [%s].", actor, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedsResponse := models.FeedsResponse{Feed: *feeds}
	respondWithJSON(w, http.StatusOK, feedsResponse)
	log.Println("Feed has been retrieved.", feeds)
}
