package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/tdanylchuk/feed-service/models"
	"github.com/tdanylchuk/feed-service/service"
	"log"
)

type KafkaFeedConsumer struct {
	FeedService service.FeedService
	Reader      *kafka.Reader
}

func CreateKafkaFeedConsumer(feedService service.FeedService, reader *kafka.Reader) *KafkaFeedConsumer {
	return &KafkaFeedConsumer{FeedService: feedService, Reader: reader}
}

func (consumer *KafkaFeedConsumer) StartConsuming() {
	log.Printf("Srarting consuming feeds from kafka...")
	go func() {
		log.Printf("Kafka feed goroutine has been started.")
		for {
			message, err := consumer.Reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("ERROR. Error while receiving message: %s", err.Error())
				continue
			}

			feedRequest := models.FeedRequest{}
			err = json.Unmarshal(message.Value, &feedRequest)
			if err != nil {
				//skipping event here, TODO: add dead letter queue
				log.Printf("ERROR. Error while unmarshalling feed: %s", err.Error())
			}

			err = consumer.FeedService.SaveFeed(feedRequest)
			if err != nil {
				//skipping event here, TODO: add dead letter queue
				log.Printf("ERROR. Error while saving feed: %s", err.Error())
			}
		}
	}()
}
