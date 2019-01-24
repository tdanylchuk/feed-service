package initializer

import (
	"github.com/segmentio/kafka-go"
	"os"
	"strings"
	"time"
)

type KafkaInitializer struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func InitKafkaClients() *KafkaInitializer {
	kafkaBrokerUrls := os.Getenv("KAFKA_HOSTS")
	feedsTopicName := os.Getenv("FEEDS_TOPIC_NAME")
	return &KafkaInitializer{
		Writer: createKafkaWriter(kafkaBrokerUrls, feedsTopicName),
		Reader: createKafkaReader(kafkaBrokerUrls, feedsTopicName),
	}
}

func createKafkaWriter(kafkaBrokerUrls string, feedsTopicName string) *kafka.Writer {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: "feed-service",
	}
	config := kafka.WriterConfig{
		Brokers:      strings.Split(kafkaBrokerUrls, ","),
		Topic:        feedsTopicName,
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	return kafka.NewWriter(config)
}

func createKafkaReader(kafkaBrokerUrls string, feedsTopicName string) *kafka.Reader {
	brokers := strings.Split(kafkaBrokerUrls, ",")
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         "feeds-service",
		Topic:           feedsTopicName,
		MinBytes:        100,
		MaxBytes:        1000,
		MaxWait:         10 * time.Microsecond,
		ReadLagInterval: -1,
	}
	return kafka.NewReader(config)
}

func (initializer *KafkaInitializer) Close() {
	defer initializer.Writer.Close()
	defer initializer.Reader.Close()
}
