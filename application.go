package main

import (
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/tdanylchuk/feed-service/controller"
	"github.com/tdanylchuk/feed-service/db"
	"github.com/tdanylchuk/feed-service/repository"
	"github.com/tdanylchuk/feed-service/service"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Router        *mux.Router
	Server        *http.Server
	DbInitializer *db.OrmPostgresInitializer
}

func CreateApp() *Application {
	dbInitializer := InitDB()
	feedController := InitController(dbInitializer.GetDB())
	router := InitRouter(feedController)

	return &Application{
		Router:        router,
		DbInitializer: dbInitializer,
	}
}

func (app *Application) Close() {
	err := app.Server.Close()
	if err != nil {
		log.Fatalln(err)
	}
	app.DbInitializer.Close()
}

func InitController(db *pg.DB) controller.FeedController {
	feedRepository := repository.CreateFeedRepository(db)
	relationRepository := repository.CreateRelationRepository(db)
	feedService := service.CreateFeedService(feedRepository, relationRepository)
	feedController := controller.CreateController(feedService)
	return feedController
}

func InitRouter(feedController controller.FeedController) *mux.Router {
	router := mux.NewRouter()
	//using {name} instead of implementing authentication
	router.
		Methods("GET").
		Path("/{actor}/feed").
		HandlerFunc(feedController.GetFeeds)
	router.
		Methods("POST").
		Path("/{actor}/feed").
		HandlerFunc(feedController.SaveFeed)
	router.
		Methods("POST").
		Path("/{actor}/action").
		HandlerFunc(feedController.PerformAction)
	router.
		Methods("GET").
		Path("/{actor}/feed/friends").
		HandlerFunc(feedController.GetFriendsFeeds)
	router.
		Methods("GET").
		Path("/{actor}/feed/friends").
		Queries("includeRelated", "{key:'^(?:tru|fals)e$}").
		HandlerFunc(feedController.GetFriendsFeeds)
	http.Handle("/", router)
	return router
}

func InitDB() *db.OrmPostgresInitializer {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_USER_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	initializer := db.New(host, user, password, dbName)
	err := initializer.InitSchema()
	if err != nil {
		panic(err)
	}
	return initializer
}

func (app *Application) StartServer(addr string) {
	log.Println("Starting feed service...")
	app.Server = &http.Server{Addr: addr, Handler: app.Router}
	log.Fatal(app.Server.ListenAndServe())
}