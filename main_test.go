package main

import (
	"context"
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app *Application
var testServer *httptest.Server

var ivansLikeFeed = map[string]string{
	"actor":  "ivan",
	"verb":   "like",
	"object": "photo:1",
	"target": "eric",
}

var ivanFollowRequest = map[string]string{
	"follow": "ivan",
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	postgresContainer := InitPostgresContainer(ctx)
	defer postgresContainer.Terminate(ctx)

	app = CreateApp()
	testServer = httptest.NewServer(app.Router)
	defer testServer.Close()
	defer app.Close()

	code := m.Run()
	os.Exit(code)
}

func TestFeedFlow(t *testing.T) {
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

func TestFollowFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//when
	expect.POST("/eric/action").
		WithJSON(ivanFollowRequest).
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
		ValueEqual("target", "ivan").
		ValueEqual("verb", "follow").
		Value("datetime").NotNull()
}

func TestForbiddenFeedPost(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/eric/feed").
		WithJSON(ivansLikeFeed).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().ValueEqual("error", "Actor[eric] is not eligible to post others[ivan] feed items.")
}

func TestFollowItself(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/ivan/action").
		WithJSON(ivanFollowRequest).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "Actor[ivan] cannot follow himself.")
}

func TestFollowEmpty(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"follow": ""}

	//then
	expect.POST("/ivan/action").
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "Follow target cannot be empty.")
}

func TestUnknownAction(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	request := map[string]string{"unknownAction": ""}

	//then
	expect.POST("/ivan/action").
		WithJSON(request).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().ValueEqual("error", "One of the actions should present. Eligible actions: 'follow'")
}

func TestBadJSONFeedPost(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/another/feed").
		WithBytes([]byte("[{{{fddfdf")).
		Expect().
		Status(http.StatusBadRequest).
		Body().Contains("Invalid request payload")
}
