package main

import (
	"log"

	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/SametAvcii/crypto-trade/pkg/server"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	consumerSuccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "consumer_success_count",
			Help: "Total number of successful Kafka consumer messages",
		},
		[]string{"consumer", "topic"},
	)
	consumerFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "consumer_failure_count",
			Help: "Total number of failed Kafka consumer messages",
		},
		[]string{"consumer", "topic"},
	)
)

func init() {
	prometheus.MustRegister(consumerSuccessCounter)
	prometheus.MustRegister(consumerFailureCounter)
}

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
	go mongoDbConsumerOrderBook.Start()
	go dbConsumerOrderBook.Start()
	//signalConsumerOrderBook.Start()
	go signalCandlesticks.Start()

	go func() {
		consumerSuccessCounter.WithLabelValues("mongoDbConsumerOrderBook", consts.OrderBookTopic).Inc()
		consumerFailureCounter.WithLabelValues("mongoDbConsumerOrderBook", consts.OrderBookTopic).Inc()
	}()

	log.Println("All consumers started successfully.")

	server.LaunchConsumerServer(config.Consumer)
}
