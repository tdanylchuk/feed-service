package sender

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/tdanylchuk/feed-service/models"
	"log"
	"time"
)

type KafkaFeedSender struct {
	Writer *kafka.Writer
}

func (sender *KafkaFeedSender) Send(request *models.FeedRequest) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Value: body,
		Time:  time.Now(),
	}
	log.Printf("Sending message to kafka - %#v.", message)
	return sender.Writer.WriteMessages(context.Background(), message)
}
