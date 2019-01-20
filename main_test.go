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

var ivansFeed = map[string]string{
	"actor":  "ivan",
	"verb":   "like",
	"object": "photo:1",
	"target": "eric",
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
		WithJSON(ivansFeed).
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

func TestForbiddenFeedPost(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)

	//then
	expect.POST("/another/feed").
		WithJSON(ivansFeed).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().ValueEqual("error", "Actor[another] is not eligible to post others[ivan] feed items.")
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
