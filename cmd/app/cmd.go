package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SametAvcii/crypto-trade/internal/clients/cache"
	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/internal/clients/kafka"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/SametAvcii/crypto-trade/pkg/server"
)

func StartApp() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	config := config.InitConfig()
	database.InitDB(config.Database)
	go database.CheckPgAlive(ctx, config.Database)

	database.InitMongo(config.Mongo)
	go database.CheckMongoAlive(ctx, config.Mongo)

	cache.InitRedis(config.Redis)
	go cache.RedisAlive(ctx, config.Redis)

	kafka.InitKafka(config.Kafka)
	go kafka.CheckKafkaAlive(ctx, config.Kafka)

	stream := events.NewStream(database.PgClient(), kafka.KafkaClientNew())

	//running all streams
	/*go func() {

		exchanges := stream.GetExchanges()
		for _, exchange := range exchanges {

			err := stream.StartAllStreams(exchange.ID.String(), consts.AggTradeTopic)
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Stream Error For Exchange: " + exchange.Name,
					Message: "Error starting stream: " + err.Error(),
					Type:    "error",
					Entity:  "exchange",
				})
				continue
			}
		}
	}()*/

	go func() {

		exchanges := stream.GetExchanges()
		for _, exchange := range exchanges {

			err := stream.StartAllStreams(exchange.ID.String(), consts.OrderBookTopic)
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Stream Error For Exchange: " + exchange.Name,
					Message: "Error starting stream: " + err.Error(),
					Type:    "error",
					Entity:  "exchange",
				})
				continue
			}
		}
	}()

	go func() {

		exchanges := stream.GetExchanges()
		for _, exchange := range exchanges {

			err := stream.StartAllStreams(exchange.ID.String(), consts.CandleStickTopic)
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Stream Error For Exchange: " + exchange.Name,
					Message: "Error starting stream: " + err.Error(),
					Type:    "error",
					Entity:  "exchange",
				})
				continue
			}
		}
	}()

	log.Println("All streams started successfully.")

	server.LaunchHttpServer(config.App, config.Allows)

	<-quit
	log.Println("Shutdown signal received. Cleaning up...")

	cancel()

}
