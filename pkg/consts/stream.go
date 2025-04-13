package consts

const (
	StreamAggTrade  = "aggTrade"
	StreamOrderBook = "depth"
)

// can define candlestick with fmt.Sprintf("%s@kline_%s", symbol, interval)
const (
	StreamCandleStick = "kline_%s"
)
