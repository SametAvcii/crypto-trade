package cmd

import (
	"log"

	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/domains/stream"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/SametAvcii/crypto-trade/pkg/server"
)

func StartApp() {
	config := config.InitConfig()
	database.InitDB(config.Database) 
	database.InitMongo(config.Mongo)
	cache.InitRedis(config.Redis)
	events.InitKafka(config.Kafka)
	stream := stream.NewStream(database.PgClient(), events.KafkaClientNew())

	go func() {

		exchanges := stream.GetExchanges()
		for _, exchange := range exchanges {

			err := stream.StartAllStreams(exchange.ID.String())
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

	go stream.Kafka.ConsumeMongoToPg()

	go stream.Kafka.ConsumeTrade()

	log.Println("All streams started successfully.")

	server.LaunchHttpServer(config.App, config.Allows)
}
