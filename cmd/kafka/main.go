package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	confluent "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	clusterID = "test-cluster"
	confluentImage = "confluentinc/confluent-local:7.5.0"
)

func main() {
	slog.Info("Spinning up a test Kafka container")
	ctx := context.Background()

	kafkaContainer, err := kafka.RunContainer(
		ctx,
		kafka.WithClusterID(clusterID),
		testcontainers.WithImage(confluentImage),
	)
	if err != nil {
		fatal("failed starting container", err)
	}
	defer terminateKafkaContainer(ctx, kafkaContainer) 

	slog.Info("Fetching test container brokers")
	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		fatal("failed fetching brokers", err)
	}
	slog.Info("Received broker values", "brokers", brokers)

	slog.Info("Configuring the Confluent config map")
	configMap := &confluent.ConfigMap{
		"bootstrap.servers": brokers[0],
	}

	slog.Info("Instantiating new admin client")
	adminClient, err := confluent.NewAdminClient(configMap)
	if err != nil {
		fatal("failed to instantiate an admin client", err)
	}

	slog.Info("Configuring topic specs")
	topicSpecs := []confluent.TopicSpecification{
		{
			Topic: "test.topic",
			NumPartitions: 1,
			ReplicationFactor: 1,
			Config: make(map[string]string),
		},
	} 

	slog.Info("Creating topics")
	topics, err := adminClient.CreateTopics(ctx, topicSpecs)
	if err != nil {
		fatal("failed to create topics", err)
	}
	for _, topic := range topics {
		slog.Info("Created topic", "name", topic.Topic)
	} 

	// Simulating some work being done
	time.Sleep(60 * time.Second)
}

func terminateKafkaContainer(ctx context.Context, container *kafka.KafkaContainer) {
	if err := container.Terminate(ctx); err != nil {
		fatal("failed terminating container", err)
	}
}

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}
