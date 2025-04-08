package events

import (
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/config"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	topics   []string
	consumer sarama.Consumer
}

var (
	kafka_client *KafkaClient
	once         sync.Once
	KafkaAlive   bool
)

func KafkaClientNew() *KafkaClient {
	return kafka_client
}
func InitKafka(cfg config.Kafka) {
	once.Do(func() {
		const (
			maxRetries    = 5
			retryInterval = 5 * time.Second
		)

		var (
			client sarama.Client
			err    error
		)

		saramaConfig := sarama.NewConfig()

		saramaConfig.Producer.Return.Errors = cfg.ReturnErrors
		saramaConfig.Producer.Return.Successes = cfg.ReturnSucces
		saramaConfig.Producer.Retry.Max = cfg.MaxRetry
		saramaConfig.Producer.MaxMessageBytes = cfg.MaxMessageSize // for large messages

		// Kafka client connection with retries
		for attempt := 1; attempt <= maxRetries; attempt++ {
			client, err = sarama.NewClient(config.ReadValue().Kafka.Brokers, saramaConfig)
			if err == nil {
				break
			}

			log.Printf("Failed to connect to Kafka client (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
		}

		if err != nil {
			log.Fatalf("Failed to connect to Kafka client after %d attempts: %v", maxRetries, err)
		}

		// Sync producer connection with retries
		var sync_producer sarama.SyncProducer
		for attempt := 1; attempt <= maxRetries; attempt++ {
			sync_producer, err = sarama.NewSyncProducerFromClient(client)
			if err == nil {
				break
			}

			log.Printf("Failed to create Kafka sync producer (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
		}

		if err != nil {
			log.Fatalf("Failed to create Kafka sync producer after %d attempts: %v", maxRetries, err)
		}

		// Get topics from Kafka with retries
		var topics []string
		for attempt := 1; attempt <= maxRetries; attempt++ {
			topics, err = client.Topics()
			if err == nil {
				break
			}

			log.Printf("Failed to get Kafka topics (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
		}

		if err != nil {
			log.Fatalf("Failed to get Kafka topics after %d attempts: %v", maxRetries, err)
		}

		// Consumer connection with retries
		var consumer sarama.Consumer
		for attempt := 1; attempt <= maxRetries; attempt++ {
			consumer, err = sarama.NewConsumerFromClient(client)
			if err == nil {
				break
			}

			log.Printf("Failed to create Kafka consumer (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
		}

		if err != nil {
			log.Fatalf("Failed to create Kafka consumer after %d attempts: %v", maxRetries, err)
		}

		// Kafka client object setup
		c := KafkaClient{
			topics:   topics,
			producer: sync_producer,
			consumer: consumer,
		}

		kafka_client = &c
		log.Println("Kafka connection initialized successfully.")
	})
}

func CheckKafkaAlive(cfg config.Kafka) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.Return.Errors = cfg.ReturnErrors
		saramaConfig.Producer.Return.Successes = cfg.ReturnSucces
		saramaConfig.Producer.Retry.Max = cfg.MaxRetry
		saramaConfig.Producer.MaxMessageBytes = cfg.MaxMessageSize

		client, err := sarama.NewClient(cfg.Brokers, saramaConfig)
		if err != nil {
			log.Println("Kafka client check failed:", err)
			KafkaAlive = false
			continue
		}

		_, err = client.Topics()
		if err != nil {
			log.Println("Kafka topic fetch failed:", err)
			KafkaAlive = false
		} else {
			KafkaAlive = true
			log.Println("Kafka is alive.")
		}

		_ = client.Close()
	}
}
