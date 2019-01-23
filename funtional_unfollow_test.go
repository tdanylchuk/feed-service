package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
	"time"
)

func AssertUnfollowFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	expect.POST("/jack/action").
		WithJSON(map[string]string{"follow": "jerry"}).
		Expect().
		Status(http.StatusOK)

	//then
	expect.POST("/jack/action").
		WithJSON(map[string]string{"unfollow": "jerry"}).
		Expect().
		Status(http.StatusOK)

	//and
	time.Sleep(asyncCallTimeout)

	//expect
	obj := expect.GET("/jack/feed").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()
	array := obj.Value("my_feed").Array()
	array.Length().Equal(2)
	array.Element(0).Object().
		ValueEqual("actor", "jack").
		NotContainsKey("object").
		ValueEqual("target", "jerry").
		ValueEqual("verb", "follow").
		Value("datetime").NotNull()
	array.Element(1).Object().
		ValueEqual("actor", "jack").
		NotContainsKey("object").
		ValueEqual("target", "jerry").
		ValueEqual("verb", "unfollow").
		Value("datetime").NotNull()
}

func AssertUnfollowItself(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/ross/action").
		WithJSON(map[string]string{"unfollow": "ross"}).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "Actor[ross] cannot unfollow himself.")
}

func AssertUnfollowEmpty(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"unfollow": ""}

	//then
	expect.POST("/ross/action").
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "unfollow target cannot be empty.")
}

func AssertUnfollowUnfollowed(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"unfollow": "someone"}

	//then
	expect.POST("/ross/action").
		WithJSON(request).
		Expect().
		Status(http.StatusInternalServerError)
}
