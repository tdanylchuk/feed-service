package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
)

var ivansLikeFeed = map[string]string{
	"actor":  "ivan",
	"verb":   "like",
	"object": "photo:1",
	"target": "eric",
}

func AssertFeedFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	obj := expect.GET("/ivan/feed").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()

	//then
	expect.POST("/ivan/feed").
		WithJSON(ivansLikeFeed).
		Expect().
		Status(http.StatusOK)

	//expect
	obj = expect.GET("/ivan/feed").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()
	array := obj.Value("my_feed").Array()
	array.Length().Equal(1)
	array.Element(0).Object().
		ValueEqual("actor", "ivan").
		ValueEqual("object", "photo:1").
		ValueEqual("target", "eric").
		ValueEqual("verb", "like").
		Value("datetime").NotNull()
}

func AssertForbiddenFeedPost(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/eric/feed").
		WithJSON(ivansLikeFeed).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().ValueEqual("error", "Actor[eric] is not eligible to post others[ivan] feed items.")
}

func AssertUnknownAction(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"unknownAction": ""}

	//then
	expect.POST("/ivan/action").
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", "One of the actions should present. Eligible actions: 'follow','unfollow'")
}

func AssertBadJSONFeedPost(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/another/feed").
		WithBytes([]byte("[{{{fddfdf")).
		Expect().
		Status(http.StatusBadRequest).
		Body().Contains("Invalid request payload")
}
