package main

import (
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App
var testServer httptest.Server

func TestMain(m *testing.M) {
	app = CreateApp()
	testServer = *httptest.NewServer(app.Router)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}

func TestFeedFlow(t *testing.T) {
	//given
	expect := httpexpect.New(t, testServer.URL)
	testFeed := map[string]string{
		"actor":  "ivan",
		"verb":   "like",
		"object": "photo:1",
		"target": "eric",
	}

	//when
	obj := expect.GET("/feed").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("my_feed").Array().Empty()
	obj.Value("next_url").String().Empty()

	//then
	expect.POST("/feed").
		WithJSON(testFeed).
		Expect().
		Status(http.StatusOK)

	//expect
	obj = expect.GET("/feed").
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