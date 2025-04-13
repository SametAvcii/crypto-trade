package main

import (
	"log"

	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/SametAvcii/crypto-trade/pkg/server"
)

func StartConsumer() {
	config := config.InitConfig()
	database.InitDB(config.Database)
	go database.CheckPgAlive(config.Database)

	database.InitMongo(config.Mongo)
	go database.CheckMongoAlive(config.Mongo)

	cache.InitRedis(config.Redis)
	go cache.RedisAlive(config.Redis)

	events.InitKafka(config.Kafka)
	go events.CheckKafkaAlive(config.Kafka)

	//stream := events.NewStream(database.PgClient(), events.KafkaClientNew())

	/*	signalConsumerAggTrade := events.Consumer{
			Brokers: config.Kafka.Brokers,
			GroupID: consts.DbOrderBookGroup,
			Topic:   consts.AggTradeTopic,
			Handler: &events.SignalHandlerAggTrade{},
		}

		dbConsumerAggTrade := events.Consumer{
			Brokers: config.Kafka.Brokers,
			GroupID: consts.DbAggTradeGroup,
			Topic:   consts.AggTradeTopic,
			Handler: &events.MongoHandler{},
		}*/

	/*signalConsumerOrderBook := events.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.SignalOrderBookGroup,
		Topic:   consts.OrderBookTopic,
		Handler: &events.SignalHandlerOrderBook{},
	}*/

	mongoDbConsumerOrderBook := events.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.DbOrderBookGroup,
		Topic:   consts.OrderBookTopic,
		Handler: &events.MongoHandler{},
	}

	dbConsumerOrderBook := events.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.DbOrderBookGroup,
		Topic:   consts.OrderBookTopic,
		Handler: &events.PgOrderBookHandler{},
	}
	signalCandlesticks := events.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.SignalCandleStickGroup,
		Topic:   consts.CandleStickTopic,
		Handler: &events.SignalHandlerCandleStick{},
	}

	// Initialize the consumer

	//dbConsumerAggTrade.Start()
	//signalConsumerAggTrade.Start()
	mongoDbConsumerOrderBook.Start()
	dbConsumerOrderBook.Start()
	//signalConsumerOrderBook.Start()
	signalCandlesticks.Start()

	log.Println("All consumers started successfully.")

	server.LaunchConsumerServer(config.Consumer)
}
