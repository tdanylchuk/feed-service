package main

import (
	"github.com/gorilla/mux"
	"github.com/tdanylchuk/feed-service/controller"
	"github.com/tdanylchuk/feed-service/repository"
	"github.com/tdanylchuk/feed-service/service"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	Server *http.Server
}

func CreateApp() App {
	app := App{}
	app.Init()
	return app
}

func (app *App) Init() {
	feedController := InitController()
	app.Router = InitRouter(feedController)
}

func InitController() controller.FeedController {
	feedRepository := repository.CreateFeedRepository()
	feedService := service.CreateFeedService(feedRepository)
	feedController := controller.CreateController(feedService)
	return feedController
}

func InitRouter(feedController controller.FeedController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feed", feedController.GetFeeds).Methods("GET")
	router.HandleFunc("/feed", feedController.SaveFeed).Methods("POST")
	http.Handle("/", router)
	return router
}

func (app *App) StartServer(addr string) {
	log.Println("Starting feed service...")
	app.Server = &http.Server{Addr: addr, Handler: app.Router}
	log.Fatal(app.Server.ListenAndServe())
}
