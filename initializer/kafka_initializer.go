package initializer

import (
	"github.com/segmentio/kafka-go"
	"log"
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
	log.Printf("Setting kafka writer to write to [%s]:[%s]", kafkaBrokerUrls, feedsTopicName)
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
		Async:        true,
	}
	return kafka.NewWriter(config)
}

func createKafkaReader(kafkaBrokerUrls string, feedsTopicName string) *kafka.Reader {
	log.Printf("Setting kafka reader to read from [%s]:[%s]", kafkaBrokerUrls, feedsTopicName)
	config := kafka.ReaderConfig{
		Brokers:         strings.Split(kafkaBrokerUrls, ","),
		GroupID:         "feeds-service",
		Topic:           feedsTopicName,
		MinBytes:        100,
		MaxBytes:        10000,
		MaxWait:         1000 * time.Second,
		ReadLagInterval: -1,
	}
	return kafka.NewReader(config)
}

func (initializer *KafkaInitializer) Close() {
	defer initializer.Writer.Close()
	defer initializer.Reader.Close()
}
