package main

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"gopkg.in/gavv/httpexpect.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app *Application
var testServer *httptest.Server

func TestMain(m *testing.M) {
	//ctx := context.Background()
	//postgresContainer := InitPostgresContainer(ctx)
	//defer postgresContainer.Terminate(ctx)
	InitPostgresRemote()

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
	testFeed := map[string]string{
		"actor":  "ivan",
		"verb":   "like",
		"object": "photo:1",
		"target": "eric",
	}

	//when
	obj := expect.GET("/feed/ivan").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	obj.Value("next_url").String().Empty()

	//then
	expect.POST("/feed").
		WithJSON(testFeed).
		Expect().
		Status(http.StatusOK)

	//expect
	obj = expect.GET("/feed/ivan").
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

func InitPostgresContainer(ctx context.Context) testcontainers.Container {
	user := "postgres_user"
	password := "postgres_password"
	dbName := "postgres_db"
	envVariables := map[string]string{
		"POSTGRES_USER":     user,
		"POSTGRES_PASSWORD": password,
		"POSTGRES_DB":       dbName,
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:9.6.8",
		ExposedPorts: []string{"5432/tcp"},
		Env:          envVariables,
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	ip, err := postgresContainer.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("%s:%s", ip, port.Port())

	_ = os.Setenv("DB_HOST", host)
	_ = os.Setenv("DB_USER", user)
	_ = os.Setenv("DB_USER_PASSWORD", password)
	_ = os.Setenv("DB_NAME", dbName)

	return postgresContainer
}

func InitPostgresRemote() {
	_ = os.Setenv("DB_HOST", "localhost:5432")
	_ = os.Setenv("DB_USER", "root")
	_ = os.Setenv("DB_USER_PASSWORD", "admin")
	_ = os.Setenv("DB_NAME", "feeds")
}
