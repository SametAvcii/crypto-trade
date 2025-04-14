package events

import (
	"testing"
	"time"

	"github.com/IBM/sarama/mocks"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInitKafka(t *testing.T) {
	// Test config
	cfg := config.Kafka{
		Brokers:        []string{"localhost:9092"},
		ReturnErrors:   true,
		ReturnSucces:   true,
		MaxRetry:       3,
		MaxMessageSize: 1000000,
	}

	// Initialize Kafka
	InitKafka(cfg)

	// Check that kafka_client was initialized
	client := KafkaClientNew()
	assert.NotNil(t, client)
	assert.NotNil(t, client.producer)
	assert.NotNil(t, client.consumer)
	assert.NotNil(t, client.topics)
}

func TestCheckKafkaAlive(t *testing.T) {
	cfg := config.Kafka{
		Brokers:        []string{"localhost:9092"},
		ReturnErrors:   true,
		ReturnSucces:   true,
		MaxRetry:       3,
		MaxMessageSize: 1000000,
	}

	// Start checking Kafka status
	go CheckKafkaAlive(cfg)

	// Give it time to run one check
	time.Sleep(time.Second)

	// Check the KafkaAlive flag
	assert.True(t, KafkaAlive)
}

func TestKafkaClientNew(t *testing.T) {
	mockClient := &KafkaClient{
		producer: mocks.NewSyncProducer(t, nil),
		consumer: mocks.NewConsumer(t, nil),
		topics:   []string{"test-topic"},
	}
	kafka_client = mockClient

	// Test KafkaClientNew
	client := KafkaClientNew()
	assert.Equal(t, mockClient, client)
	assert.NotNil(t, client.producer)
	assert.NotNil(t, client.consumer)
	assert.Equal(t, []string{"test-topic"}, client.topics)
}
