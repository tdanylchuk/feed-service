package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"time"
)

func InitKafkaContainers(ctx context.Context) (testcontainers.Container, testcontainers.Container) {
	zookeeperContainer, host := InitZookeeperContainer(ctx)
	kafkaContainer := InitKafkaContainer(ctx, host)
	return zookeeperContainer, kafkaContainer
}

func InitZookeeperContainer(ctx context.Context) (testcontainers.Container, string) {
	log.Printf("Starting Zookeeper container...")
	req := testcontainers.ContainerRequest{
		Image:        "wurstmeister/zookeeper",
		ExposedPorts: []string{"2181:2181"},
		Env:          map[string]string{},
	}
	container, host := StartContainerAndGetHost(ctx, req, "2181")
	_ = os.Setenv("ZOOKEEPER_HOSTS", host)
	log.Printf("Test container with Zookeeper has been started. Zookeeper host[%s].", host)
	return container, host
}

func InitKafkaContainer(ctx context.Context, zookeeperHost string) testcontainers.Container {
	log.Printf("Starting Kafka container...")
	feedsTopic := "feeds"
	envVariables := map[string]string{
		"KAFKA_CREATE_TOPICS":                    fmt.Sprintf("%s:1:1", feedsTopic),
		"KAFKA_ZOOKEEPER_CONNECT":                zookeeperHost,
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":   "INTERNAL_PLAINTEXT:PLAINTEXT",
		"KAFKA_ADVERTISED_LISTENERS":             "INTERNAL_PLAINTEXT://192.168.99.100:9092",
		"KAFKA_LISTENERS":                        "INTERNAL_PLAINTEXT://0.0.0.0:9092",
		"KAFKA_INTER_BROKER_LISTENER_NAME":       "INTERNAL_PLAINTEXT",
		"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
		"KAFKA_OFFSETS_TOPIC_NUM_PARTITIONS":     "1",
		"KAFKA_LOG_FLUSH_INTERVAL_MS":            "1",
		"KAFKA_REPLICA_SOCKET_TIMEOUT_MS":        "1000",
		"KAFKA_CONTROLLER_SOCKET_TIMEOUT_MS":     "1000",
	}
	req := testcontainers.ContainerRequest{
		Image:        "wurstmeister/kafka:2.12-2.1.0",
		ExposedPorts: []string{"9092:9092"},
		Env:          envVariables,
		WaitingFor:   &KafkaWaitStrategy{TopicToTest: feedsTopic},
	}
	container, host := StartContainerAndGetHost(ctx, req, "9092")

	_ = os.Setenv("KAFKA_HOSTS", host)
	_ = os.Setenv("FEEDS_TOPIC_NAME", feedsTopic)

	log.Printf("Test container with Kafka has been started. Kafka host[%s].", host)
	return container
}

type KafkaWaitStrategy struct {
	TopicToTest string
}

func (strategy *KafkaWaitStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	ip, _ := target.Host(ctx)
	port, _ := target.MappedPort(ctx, "9092")
	host := fmt.Sprintf("%s:%s", ip, port.Port())

	attempt := 0
	for attempt < 60 {
		attempt++
		log.Printf("Checking Kafka Up&Running on host [%s]... Attempt #[%d]", host, attempt)
		err := strategy.CheckKafkaConnectionAndCreateTopic(ctx, host)
		if err == nil {
			log.Println("Successfully connected to Kafka.")
			return err
		}
		log.Println(err)
		time.Sleep(time.Second)
	}
	return errors.New(fmt.Sprintf("failed to connect to Kafka container using host[%s] ", host))
}

func (strategy *KafkaWaitStrategy) CheckKafkaConnectionAndCreateTopic(ctx context.Context, host string) error {
	partitions, err := kafka.LookupPartitions(ctx, "tcp", host, strategy.TopicToTest)
	if err != nil {
		return err
	}
	if len(partitions) < 1 {
		return errors.New(fmt.Sprintf(
			"At least a single partition should be assigned to test topic [%s]", strategy.TopicToTest))
	}
	return nil
}
