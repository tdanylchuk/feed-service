package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	var feed models.Feed
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&feed); err != nil {
		log.Fatalln("Something went wrong during saving feed...", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := controller.FeedService.SaveFeed(feed); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (controller *DefaultFeedController) GetFeeds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name, _ := vars["name"]
	log.Printf("Controller. Retrieving user feed for [%s]...", name)
	feeds, err := controller.FeedService.RetrieveFeed(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedsResponse := models.FeedsResponse{Feed: *feeds}
	respondWithJSON(w, http.StatusOK, feedsResponse)
	log.Println("Feed has been retrieved.", feeds)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
