package consts

const (
	// KafkaTopicTrade is the topic name for trade data
	AggTradeTopic    = "agg-trade-data"
	OrderBookTopic   = "depth-data"
	CandleStickTopic = "candlestick-data"
)

// pg topic
const (
	PgAggTradeTopic    = "agg-trade-data-pg"
	PgOrderBookTopic   = "depth-data-pg"
	PgCandleStickTopic = "candlestick-data-pg"
)

const ( // Signal Group
	SignalCandleStickGroup = "signal-candle-stick-group"
	SignalOrderBookGroup   = "signal-order-book-group"
	MongoCandleStickGroup  = "mongo-candle-stick-group"
)

const ( // DB Group
	DbCandleStickGroup = "db-candle-stick-group"
	DbOrderBookGroup   = "db-order-book-group"
	PgOrderBookGroup   = "pg-order-book-group"
	PgCandleStickGroup = "pg-candlestick-group"
)
