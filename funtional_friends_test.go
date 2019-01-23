package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
	"time"
)

func AssertRetrieveFriendsFeed(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	tarasLikeFeed := map[string]string{
		"actor":  "taras",
		"verb":   "like",
		"object": "photo:2",
		"target": "eric",
	}
	igorsShareFeed := map[string]string{
		"actor":  "igor",
		"verb":   "share",
		"object": "post:2",
	}

	//when
	expect.POST("/taras/feed").
		WithJSON(tarasLikeFeed).
		Expect().
		Status(http.StatusOK)
	//and
	expect.POST("/igor/action").
		WithJSON(map[string]string{"follow": "ivan"}).
		Expect().
		Status(http.StatusOK)
	//and
	expect.POST("/igor/feed").
		WithJSON(igorsShareFeed).
		Expect().
		Status(http.StatusOK)

	//then
	expect.POST("/andrew/action").
		WithJSON(map[string]string{"follow": "taras"}).
		Expect().
		Status(http.StatusOK)
	//and
	time.Sleep(asyncCallTimeout)
	//expect
	friendsFeed := expect.GET("/andrew/feed/friends").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	friendsFeed.Value("next_url").String().Empty()
	friendsFeed.Value("friends_feed").Array().
		Length().Equal(1)

	//then
	expect.POST("/andrew/action").
		WithJSON(map[string]string{"follow": "igor"}).
		Expect().
		Status(http.StatusOK)
	//and
	time.Sleep(asyncCallTimeout)
	//expect
	friendsFeed = expect.GET("/andrew/feed/friends").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	friendsFeed.Value("next_url").String().Empty()
	friendsFeed.Value("friends_feed").Array().
		Length().Equal(3)

	//then
	expect.POST("/andrew/action").
		WithJSON(map[string]string{"unfollow": "taras"}).
		Expect().
		Status(http.StatusOK)
	//and
	time.Sleep(asyncCallTimeout)
	//expect
	friendsFeed = expect.GET("/andrew/feed/friends").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	friendsFeed.Value("next_url").String().Empty()
	friendsFeed.Value("friends_feed").Array().
		Length().Equal(2)
}
