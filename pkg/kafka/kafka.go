package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/config"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	topics   []string
}

var (
	kafka_client *KafkaClient
	once         sync.Once
)

func KafkaClientNew() *KafkaClient {
	return kafka_client
}

func InitKafka(kafka_config config.Kafka) {
	once.Do(func() {
		var (
			client sarama.Client
			err    error
		)

		saramaConfig := sarama.NewConfig()

		saramaConfig.Producer.Return.Errors = kafka_config.ReturnErrors
		saramaConfig.Producer.Return.Successes = kafka_config.ReturnSucces
		saramaConfig.Producer.Retry.Max = kafka_config.MaxRetry
		saramaConfig.Producer.MaxMessageBytes = kafka_config.MaxMessageSize // for large messages

		client, err = sarama.NewClient(config.ReadValue().Kafka.Brokers, saramaConfig)
		if err != nil {
			log.Fatalf("Could not connect to kafka: %v", err)
		}

		sync_producer, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			log.Fatalf("Could not connect to kafka: %v", err)
		}

		topics, err := client.Topics()
		if err != nil {
			log.Fatalf("Could not connect to kafka: %v", err)
		}

		c := KafkaClient{
			topics:   topics,
			producer: sync_producer,
		}

		kafka_client = &c
	})
}

func (k *KafkaClient) Produce(topic, key string, message []byte) (int32, int64, error) {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	return k.producer.SendMessage(msg)
}
