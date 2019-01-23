package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/tdanylchuk/feed-service/consumer"
	"github.com/tdanylchuk/feed-service/controller"
	"github.com/tdanylchuk/feed-service/initializer"
	"github.com/tdanylchuk/feed-service/repository"
	"github.com/tdanylchuk/feed-service/sender"
	"github.com/tdanylchuk/feed-service/service"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Router            *mux.Router
	Server            *http.Server
	DbInitializer     *initializer.OrmPostgresInitializer
	KafkaInitializer  *initializer.KafkaInitializer
	KafkaFeedConsumer *consumer.KafkaFeedConsumer
}

func CreateApp() *Application {
	dbInitializer := initializer.InitPostgresDB()
	kafkaInitializer := initializer.InitKafkaClients()
	feedService := InitFeedService(dbInitializer.GetDB(), kafkaInitializer.Writer)
	kafkaFeedConsumer := InitKafkaFeedConsumer(feedService, kafkaInitializer.Reader)
	feedController := InitController(feedService)
	router := InitRouter(feedController)

	return &Application{
		Router:            router,
		DbInitializer:     dbInitializer,
		KafkaInitializer:  kafkaInitializer,
		KafkaFeedConsumer: kafkaFeedConsumer,
	}
}

func (app *Application) Close() {
	defer app.Server.Close()
	defer app.DbInitializer.Close()
	defer app.KafkaInitializer.Close()
}

func InitFeedService(db *pg.DB, kafkaWriter *kafka.Writer) service.FeedService {
	feedRepository := repository.CreateFeedRepository(db)
	relationRepository := repository.CreateRelationRepository(db)
	eventSender := sender.CreateKafkaSender(kafkaWriter)
	return service.CreateFeedService(feedRepository, relationRepository, eventSender)
}

func InitController(feedService service.FeedService) controller.FeedController {
	return controller.CreateController(feedService)
}

func InitKafkaFeedConsumer(feedService service.FeedService, reader *kafka.Reader) *consumer.KafkaFeedConsumer {
	return consumer.CreateKafkaFeedConsumer(feedService, reader)
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
		HandlerFunc(feedController.ProcessFeed)
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

func (app *Application) StartServer() {
	port, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		port = "8000"
	}
	log.Printf("Starting feed service on port [%s]...", port)
	app.Server = &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: app.Router}
	app.KafkaFeedConsumer.StartConsuming()
	log.Fatal(app.Server.ListenAndServe())
}
