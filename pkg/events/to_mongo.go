package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoHandler struct{}

func (d *MongoHandler) HandleMessage(msg *sarama.ConsumerMessage) {
	//log.Printf("[Mongo] Processing: %s", msg.Topic)
	mongoData, err := insertMessageToMongo(msg)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting message into MongoDB",
			Message: "Error inserting message into MongoDB: " + err.Error(),
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})
		log.Printf("Error inserting message into MongoDB: %v", err)
		return
	}
	mongoID := mongoData.InsertedID.(primitive.ObjectID).Hex()
	aggTradeMongo := dtos.MongoData{
		MongoID: mongoID,
		Value:   string(msg.Value),
	}
	jsonBytes, err := json.Marshal(aggTradeMongo)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error marshalling Kafka message",
			Message: "Error marshalling Kafka message: " + err.Error(),
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})

		log.Printf("Error marshaling Kafka message: %v", err)
		return
	}

	client := KafkaClientNew()
	pgTopic := getPgTopic(msg.Topic)
	if pgTopic == "" {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error getting Postgres topic",
			Message: "Error getting Postgres topic for Kafka message: ",
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})
		log.Printf("Error getting Postgres topic for Kafka message: %v", err)
		return
	}

	_, _, err = client.Produce(pgTopic, mongoID, jsonBytes)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error sending message to Kafka",
			Message: "Error sending message to Kafka: " + err.Error(),
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})
		log.Printf("Error sending message to postgres-topic: %v", err)
	} else {
		//log.Printf("Message sent to postgres-topic at partition %d offset %d", partition, offset)
	}
	//log.Printf("MongoDB updated for id %s", mongoID)
}

func insertMessageToMongo(msg *sarama.ConsumerMessage) (*mongo.InsertOneResult, error) {
	var doc bson.M
	if err := bson.UnmarshalExtJSON(msg.Value, true, &doc); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling message",
			Message: "Error unmarshalling message into BSON: " + err.Error(),
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})
		log.Printf("Error unmarshalling message into BSON: %v", err)
		return nil, err
	}

	mongoClient := database.MongoClient()
	collName := getCollectionName(msg.Topic)
	collection := mongoClient.Database(config.ReadValue().Mongo.Database).Collection(collName)

	mongoData, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting message into MongoDB",
			Message: "Error inserting message into MongoDB: " + err.Error(),
			Type:    "error",
			Entity:  "trade",
			Data:    string(msg.Value),
		})
		log.Printf("Error inserting message into MongoDB: %v", err)
		return nil, err
	}

	return mongoData, nil
}

func getPgTopic(topic string) string {
	switch topic {
	case consts.AggTradeTopic:
		return consts.PgAggTradeTopic
	case consts.OrderBookTopic:
		return consts.PgOrderBookTopic
	default:
		return ""
	}
}

func getCollectionName(topic string) string {
	switch topic {
	case consts.AggTradeTopic:
		return consts.CollectionNameTrade
	case consts.OrderBookTopic:
		return consts.CollectionNameOrder
	default:
		return ""
	}
}
