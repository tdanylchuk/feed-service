package main

import (
	"fmt"
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"testing"
	"time"
)

func AssertFeedPaginationFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	expect.POST("/alex/feed").
		WithJSON(getFeed("photo:10")).
		Expect().
		Status(http.StatusOK)
	expect.POST("/alex/feed").
		WithJSON(getFeed("photo:11")).
		Expect().
		Status(http.StatusOK)
	expect.POST("/alex/feed").
		WithJSON(getFeed("photo:12")).
		Expect().
		Status(http.StatusOK)
	//and
	time.Sleep(asyncCallTimeout * 3)

	//expect
	obj := expect.GET("/alex/feed").WithQuery("page", 1).WithQuery("limit", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("my_feed").Array().Length().Equal(2)
	nextUrlRaw := obj.Value("next_url").String().
		NotEmpty().
		Equal(fmt.Sprintf("%s/alex/feed?limit=2&page=2", testServer.URL)).
		Raw()

	//and
	obj = expect.GET("").WithURL(nextUrlRaw).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().NotEmpty().Equal(fmt.Sprintf("%s/alex/feed?limit=2&page=3", testServer.URL))
	obj.Value("my_feed").Array().Length().Equal(1)
}

func getFeed(object string) map[string]string {
	return map[string]string{
		"actor":  "alex",
		"verb":   "like",
		"object": object,
	}
}
