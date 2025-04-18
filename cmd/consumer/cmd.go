package main

import (
	"log"

	"github.com/SametAvcii/crypto-trade/internal/clients/cache"
	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/internal/clients/kafka"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
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

	kafka.InitKafka(config.Kafka)
	go kafka.CheckKafkaAlive(config.Kafka)

	//TODO: fix this tomorrow
	mongoDbConsumerOrderBook := kafka.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.DbOrderBookGroup,
		Topic:   consts.OrderBookTopic,
		Handler: &events.MongoHandler{},
	}

	dbConsumerOrderBook := kafka.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.PgOrderBookGroup,
		Topic:   consts.PgOrderBookTopic,
		Handler: &events.PgOrderBookHandler{},
	}

	mongoDbConsumerCandleStick := kafka.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.MongoCandleStickGroup,
		Topic:   consts.CandleStickTopic,
		Handler: &events.MongoHandler{},
	}

	dbConsumerCandlestick := kafka.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.PgCandleStickGroup,
		Topic:   consts.PgCandleStickTopic,
		Handler: &events.PgCandleStickHandler{},
	}

	signalCandlesticks := kafka.Consumer{
		Brokers: config.Kafka.Brokers,
		GroupID: consts.SignalCandleStickGroup,
		Topic:   consts.CandleStickTopic,
		Handler: &events.SignalHandlerCandleStick{},
	}

	// Initialize the consumer

	mongoDbConsumerOrderBook.Start()
	dbConsumerOrderBook.Start()

	mongoDbConsumerCandleStick.Start()
	dbConsumerCandlestick.Start()
	signalCandlesticks.Start()

	// Start the consumers

	go func() {
		consumerSuccessCounter.WithLabelValues("mongoDbConsumerOrderBook", consts.OrderBookTopic).Inc()
		consumerFailureCounter.WithLabelValues("mongoDbConsumerOrderBook", consts.OrderBookTopic).Inc()
	}()

	go func() {
		consumerSuccessCounter.WithLabelValues("dbConsumerOrderBook", consts.OrderBookTopic).Inc()
		consumerFailureCounter.WithLabelValues("dbConsumerOrderBook", consts.OrderBookTopic).Inc()
	}()

	go func() {
		consumerSuccessCounter.WithLabelValues("signalCandlesticks", consts.CandleStickTopic).Inc()
		consumerFailureCounter.WithLabelValues("signalCandlesticks", consts.CandleStickTopic).Inc()
	}()

	log.Println("All consumers started successfully.")

	server.LaunchConsumerServer(config.Consumer)
}
