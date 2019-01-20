package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
)

var rossFollowRequest = map[string]string{
	"follow": "ross",
}

func AssertFollowFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	expect.POST("/eric/action").
		WithJSON(rossFollowRequest).
		Expect().
		Status(http.StatusOK)

	//then
	obj := expect.GET("/eric/feed").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()
	array := obj.Value("my_feed").Array()
	array.Length().Equal(1)
	array.Element(0).Object().
		ValueEqual("actor", "eric").
		NotContainsKey("object").
		ValueEqual("target", "ross").
		ValueEqual("verb", "follow").
		Value("datetime").NotNull()
}

func AssertFollowItself(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/ross/action").
		WithJSON(rossFollowRequest).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "Actor[ross] cannot follow himself.")
}

func AssertFollowEmpty(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"follow": ""}

	//then
	expect.POST("/ivan/action").
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "follow target cannot be empty.")
}

func AssertDoubleFollow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	expect.POST("/wolf/action").
		WithJSON(rossFollowRequest).
		Expect().
		Status(http.StatusOK)

	//then
	expect.POST("/wolf/action").
		WithJSON(rossFollowRequest).
		Expect().
		Status(http.StatusInternalServerError)
}
