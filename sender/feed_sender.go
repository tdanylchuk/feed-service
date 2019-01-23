package sender

import (
	"github.com/segmentio/kafka-go"
	"github.com/tdanylchuk/feed-service/models"
)

type Sender interface {
	Send(request *models.FeedRequest) error
}

func CreateKafkaSender(writer *kafka.Writer) Sender {
	return &KafkaFeedSender{Writer: writer}
}
