package main

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var app *Application
var testServer *httptest.Server

//in case of instability of tests, please tweak these timeouts
var serverWarmupTimeout = 10 * time.Second
var asyncCallTimeout = 1000 * time.Millisecond

func TestMain(m *testing.M) {
	ctx := context.Background()
	postgresContainer := InitPostgresContainer(ctx)
	zookeeperContainer, kafkaContainer := InitKafkaContainers(ctx)
	defer postgresContainer.Terminate(ctx)
	defer zookeeperContainer.Terminate(ctx)
	defer kafkaContainer.Terminate(ctx)

	app = CreateApp()
	app.KafkaFeedConsumer.StartConsuming()

	testServer = httptest.NewServer(app.Router)
	defer testServer.Close()
	defer app.Close()

	time.Sleep(serverWarmupTimeout)
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

//feed pagination suit
func TestFeedPaginationFlow(t *testing.T) {
	AssertFeedPaginationFlow(t)
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
