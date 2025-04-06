package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	topics   []string
	consumer sarama.Consumer
}

var (
	kafka_client *KafkaClient
	once         sync.Once
)

func KafkaClientNew() *KafkaClient {
	return kafka_client
}

func InitKafka(cfg config.Kafka) {
	once.Do(func() {
		var (
			client sarama.Client
			err    error
		)

		saramaConfig := sarama.NewConfig()

		saramaConfig.Producer.Return.Errors = cfg.ReturnErrors
		saramaConfig.Producer.Return.Successes = cfg.ReturnSucces
		saramaConfig.Producer.Retry.Max = cfg.MaxRetry
		saramaConfig.Producer.MaxMessageBytes = cfg.MaxMessageSize // for large messages

		client, err = sarama.NewClient(config.ReadValue().Kafka.Brokers, saramaConfig)
		if err != nil {
			log.Fatalf("Could not connect to kafka client: %v", err)
		}

		sync_producer, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			log.Fatalf("Could not connect to kafka producer: %v", err)
		}

		topics, err := client.Topics()
		if err != nil {
			log.Fatalf("Could not connect to kafka topics: %v", err)
		}

		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			log.Fatalf("Could not connect to kafka consumer: %v", err)
		}

		c := KafkaClient{
			topics:   topics,
			producer: sync_producer,
			consumer: consumer,
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

func (k *KafkaClient) ConsumeTrade() error {

	topic := "crypto-trade"
	log.Println("Starting Kafka consumer...", topic)
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("Error creating Kafka consumer: %v", err)
	}
	defer partitionConsumer.Close()

	//kafka policy
	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))
		var doc bson.M
		if err := bson.UnmarshalExtJSON(msg.Value, true, &doc); err != nil {
			log.Printf("Error unmarshalling message into BSON: %v", err)
			continue
		}
		mongoClient := database.MongoClient()
		collection := mongoClient.Database("crypto-trade").Collection("crypto-trade")
		mongoData, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			log.Printf("Error inserting message into MongoDB: %v", err)

		} else {

			mongoID := mongoData.InsertedID.(primitive.ObjectID).Hex()

			aggTradeMongo := dtos.AggTradeMongo{
				MongoID: mongoID,
				Value:   string(msg.Value),
			}

			jsonBytes, err := json.Marshal(aggTradeMongo)
			if err != nil {
				log.Printf("Error marshaling Kafka message: %v", err)
				continue
			}

			partition, offset, err := k.Produce("to-postgres-price", mongoID, jsonBytes)
			if err != nil {
				log.Printf("Error sending message to postgres-topic: %v", err)
			} else {
				log.Printf("Message sent to postgres-topic at partition %d offset %d", partition, offset)
			}

			log.Printf("MongoDB updated for id %s with value %s", mongoID, string(msg.Value))

		}
	}

	return nil
}

func (k *KafkaClient) ConsumeMongoToPg() error {
	topic := "to-postgres-price"

	log.Println("Starting Kafka consumer...", topic)

	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("Error creating Kafka consumer: %v", err)
	}
	log.Printf("Received message for Postgres: %s", topic)
	defer partitionConsumer.Close()
	log.Printf("Starting Kafka consumer for Postgres...")
	var (
		payload     dtos.AggTradeMongo
		symbolPrice entities.SymbolPrice
		aggTrade    dtos.AggTrade
	)

	for msg := range partitionConsumer.Messages() {

		log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Error unmarshalling message for Postgres: %v", err)

		}

		db := database.PgClient()

		json.Unmarshal([]byte(payload.Value), &aggTrade)
		symbolPrice.FromDto(&aggTrade)
		symbolPrice.MongoID = payload.MongoID

		if err := db.Create(&symbolPrice).Error; err != nil {
			log.Printf("Error inserting message into Postgres: %v", err)

		} else {
			log.Printf("Postgres updated for id %d with value %s", symbolPrice.ID, payload.Value)
		}
	}
	return nil
}
