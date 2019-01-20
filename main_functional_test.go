package main

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
)

var app *Application
var testServer *httptest.Server

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

//feed suit
func TestFeedFlow(t *testing.T) {
	AssertFeedFlow(t)
}
func TestForbiddenFeedPost(t *testing.T) {
	AssertForbiddenFeedPost(t)
}
func TestUnknownAction(t *testing.T) {
	AssertUnknownAction(t)
}
func TestBadJSONFeedPost(t *testing.T) {
	AssertBadJSONFeedPost(t)
}

//enriched feed suit
func TestEnrichedFeedFlow(t *testing.T) {
	AssertEnrichedFeedFlow(t)
}

//friends feed suit
func TestRetrieveFriendsFeed(t *testing.T) {
	AssertRetrieveFriendsFeed(t)
}

//follow test suit
func TestFollowFlow(t *testing.T) {
	AssertFollowFlow(t)
}
func TestFollowItself(t *testing.T) {
	AssertFollowItself(t)
}
func TestFollowEmpty(t *testing.T) {
	AssertFollowEmpty(t)
}
func TestDoubleFollow(t *testing.T) {
	AssertDoubleFollow(t)
}

//unfollow test suit
func TestUnfollowFlow(t *testing.T) {
	AssertUnfollowFlow(t)
}
func TestUnfollowItself(t *testing.T) {
	AssertUnfollowItself(t)
}
func TestUnfollowEmpty(t *testing.T) {
	AssertUnfollowEmpty(t)
}
func TestUnfollowUnfollowed(t *testing.T) {
	AssertUnfollowUnfollowed(t)
}
