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

func (controller *DefaultFeedController) ProcessFeed(w http.ResponseWriter, r *http.Request) {
	log.Println("Controller. Processing new feed...")
	actor := getActor(r)
	var feed models.FeedRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&feed); err != nil {
		log.Println("Something went wrong during decoding feed...", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if feed.Actor != actor {
		str := fmt.Sprintf("Actor[%s] is not eligible to post others[%s] feed items.", actor, feed.Actor)
		respondWithError(w, http.StatusForbidden, str)
		return
	}

	if err := controller.FeedService.ProcessFeed(feed); err != nil {
		log.Printf("Something went wrong during feed processing. Feed - [%s]. Error - [%s]", feed, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (controller *DefaultFeedController) GetFeeds(w http.ResponseWriter, r *http.Request) {
	actor := getActor(r)
	includeRelated := getBoolParam(r, "includeRelated")
	log.Printf("Controller. Retrieving user feed for [%s] with related included[%t]...", actor, includeRelated)

	page, limit, err := getPagingValues(r)
	if err != nil {
		str := fmt.Sprintf("Incorrect pagination values in request")
		respondWithError(w, http.StatusBadRequest, str)
	}

	feeds, err := controller.FeedService.RetrieveFeed(actor, includeRelated, page, limit)
	if err != nil {
		log.Printf("Something went wrong during retriving feeds for [%s]. Error - [%s].", actor, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nextUrl := getNextUrl(r, page, limit)
	feedsResponse := models.FeedsResponse{Feed: *feeds, NextUrl: nextUrl}
	respondWithJSON(w, http.StatusOK, feedsResponse)
	log.Println("Feed has been retrieved.", feedsResponse)
}

func (controller *DefaultFeedController) PerformAction(w http.ResponseWriter, r *http.Request) {
	actor := getActor(r)
	log.Printf("Controller. Processing action request from [%s]...", actor)

	var actionRequest models.ActionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&actionRequest); err != nil {
		log.Println("Something went wrong during decoding action request...", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := actionRequest.Validate(actor); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := controller.FeedService.ProcessAction(actor, actionRequest); err != nil {
		str := fmt.Sprintf("Something went wrong during processing action request. Request - [%#v]. Error - [%s]",
			actionRequest, err)
		respondWithError(w, http.StatusInternalServerError, str)
		return
	}
}

func (controller *DefaultFeedController) GetFriendsFeeds(w http.ResponseWriter, r *http.Request) {
	actor := getActor(r)
	log.Printf("Controller. Retrieving friends feed for [%s]...", actor)

	page, limit, err := getPagingValues(r)
	if err != nil {
		str := fmt.Sprintf("Incorrect pagination values in request")
		respondWithError(w, http.StatusBadRequest, str)
	}

	feeds, err := controller.FeedService.RetrieveFriendsFeed(actor, page, limit)
	if err != nil {
		log.Printf("Something went wrong during retriving friends feeds for [%s]. Error - [%s].", actor, err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nextUrl := getNextUrl(r, page, limit)
	feedsResponse := models.FriendsFeedsResponse{Feed: *feeds, NextUrl: nextUrl}
	respondWithJSON(w, http.StatusOK, feedsResponse)
	log.Println("Friends feed has been retrieved.", feedsResponse)
}
