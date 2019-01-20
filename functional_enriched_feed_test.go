package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
)

func AssertEnrichedFeedFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	object := "photo:1"

	//when
	expect.POST("/scott/feed").
		WithJSON(map[string]string{
			"actor":  "scott",
			"verb":   "post",
			"object": object}).
		Expect().
		Status(http.StatusOK)
	//and
	expect.POST("/smith/feed").
		WithJSON(map[string]string{
			"actor":  "smith",
			"verb":   "share",
			"object": object,
			"target": "someone"}).
		Expect().
		Status(http.StatusOK)

	//then
	expect.POST("/aron/feed").
		WithJSON(map[string]string{
			"actor":  "aron",
			"verb":   "like",
			"object": object,
			"target": "someone"}).
		Expect().
		Status(http.StatusOK)

	//expect
	expect.GET("/aron/feed").WithQuery("includeRelated", "true").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		Value("my_feed").Array().
		Element(0).Object().
		NotContainsKey("related")

	//then
	expect.POST("/aron/action").
		WithJSON(map[string]string{
			"follow": "scott"}).
		Expect().
		Status(http.StatusOK)

	//expect
	obj := expect.GET("/aron/feed").WithQuery("includeRelated", "true").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()
	array := obj.Value("my_feed").Array()
	array.Length().Equal(2)
	feed := array.Element(0).Object().
		ValueEqual("actor", "aron").
		ValueEqual("object", object).
		ValueEqual("target", "someone").
		ValueEqual("verb", "like")
	feed.Value("datetime").NotNull()
	relatedArray := feed.Value("related").Array()
	relatedArray.Length().Equal(1)
	relatedArray.Element(0).Object().
		ValueEqual("actor", "scott").
		ValueEqual("object", object).
		ValueEqual("verb", "post").
		NotContainsKey("target")
}
