package main

import (
	"log"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/SametAvcii/crypto-trade/pkg/server"
)

func StartConsumer() {
	config := config.InitConfig()
	database.InitDB(config.Database)
	database.InitMongo(config.Mongo)
	events.InitKafka(config.Kafka)

	stream := events.NewStream(database.PgClient(), events.KafkaClientNew())

	// Initialize the consumer
	go func() {
		if err := stream.Kafka.ConsumeMongoToPg(); err != nil {
			log.Fatalf("Error consuming MongoDB to PostgreSQL: %v", err)
		}
	}()

	go func() {
		if err := stream.Kafka.ConsumeTrade(); err != nil {
			log.Fatalf("Error consuming trade data: %v", err)
		}
	}()

	log.Println("All consumers started successfully.")

	server.LaunchConsumerServer(config.Consumer)
}
