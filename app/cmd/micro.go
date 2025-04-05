package cmd

import (
	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/domains/stream"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/SametAvcii/crypto-trade/pkg/kafka"
	"github.com/SametAvcii/crypto-trade/pkg/server"
)

func StartApp() {
	config := config.InitConfig()
	database.InitDB(config.Database)
	cache.InitRedis(config.Redis)
	kafka.InitKafka(config.Kafka)

	go func() {
		stream := stream.NewStream(database.PgClient())

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

	server.LaunchHttpServer(config.App, config.Allows)
}
