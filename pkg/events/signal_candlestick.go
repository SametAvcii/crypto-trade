package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/candlestick"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type SignalHandlerCandleStick struct{}

func (s *SignalHandlerCandleStick) HandleMessage(msg *sarama.ConsumerMessage) {

	var payload dtos.CandlestickWs
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error unmarshalling message for candlestick signal",
			Message: "Error unmarshalling message into BSON: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(msg.Value),
		})
		return
	}

	// Check if the candlestick is closed
	if payload.Kline.IsKlineClosed == false {
		log.Println("Candlestick is not closed, skipping...")
		return
	}

	var pgDb = database.PgClient()
	var interval entities.SignalInterval

	err := pgDb.Where("symbol = ? and interval = ?", strings.ToLower(payload.Symbol), payload.Kline.Interval).First(&interval).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching intervals from Postgres",
			Message: "Error fetching intervals from Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "aggTrade",
			Data:    string(msg.Value),
		})
		return
	}

	key50 := fmt.Sprintf("%s:%s:ma50", payload.Symbol, interval.Interval)
	key200 := fmt.Sprintf("%s:%s:ma200", payload.Symbol, interval.Interval)

	rdb := cache.RedisClient()
	pipe := rdb.Pipeline()
	price, _ := decimal.NewFromString(payload.Kline.ClosePrice)

	if exists, _ := rdb.Exists(context.Background(), key50).Result(); exists == 0 {
		if exists, _ := rdb.Exists(context.Background(), key200).Result(); exists > 0 {
			pipe.Del(context.Background(), key200)
		}

		var candleSticks []entities.Candlestick
		err = pgDb.Where("symbol = ? and interval = ?", payload.Symbol, payload.Kline.Interval).
			Order("open_time desc").Limit(200).Find(&candleSticks).Error

		if err != nil || len(candleSticks) < 200 {
			ctlog.CreateLog(&entities.Log{
				Title:   "Error fetching candlesticks from Postgres",
				Message: "Error fetching candlesticks from Postgres:",
				Type:    "error",
				Entity:  "candlestick",
				Data:    string(msg.Value),
			})

			candleSticks, err = candlestick.GetCandleSticksAndUpdate(interval.ExchangeID.String(), payload.Symbol, interval.Interval, 200)
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Error fetching candlesticks from API",
					Message: "Error fetching candlesticks from API: " + err.Error(),
					Type:    "error",
					Entity:  "candlestick",
					Data:    string(msg.Value),
				})
				return
			}

		}

		var cnt = 0
		for _, candle := range candleSticks {
			if cnt < 50 {
				pipe.RPush(context.Background(), key50, candle.Close.String())
				cnt++
			}
			pipe.RPush(context.Background(), key200, candle.Close.String())
		}

	} else {
		pipe.RPush(context.Background(), key50, price.String())
		pipe.LTrim(context.Background(), key50, -50, -1)

		pipe.RPush(context.Background(), key200, price.String())
		pipe.LTrim(context.Background(), key200, -200, -1)
	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error executing Redis pipeline",
			Message: "Error executing Redis pipeline: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    string(msg.Value),
		})
	}

	signal, err := checkForSignal(rdb, payload.Symbol, interval.Interval)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error checking for signal",
			Message: "Error checking for signal: " + err.Error(),
			Type:    "error",
			Entity:  "signal",
			Data:    string(msg.Value),
		})
		return
	}
	log.Printf("[%s][%s] Signal: %s", payload.Symbol, interval.Interval, signal)

}

func checkForSignal(rdb *redis.Client, symbol, timeframe string) (dtos.MaData, error) {
	var res dtos.MaData
	key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
	key200 := fmt.Sprintf("%s:%s:ma200", symbol, timeframe)
	lastSignalKey := fmt.Sprintf("%s:%s:lastSignal", symbol, timeframe)

	pipe := rdb.Pipeline()
	res50 := pipe.LRange(context.Background(), key50, 0, -1)
	res200 := pipe.LRange(context.Background(), key200, 0, -1)

	_, err := pipe.Exec(context.Background())
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error executing Redis pipeline",
			Message: "Error executing Redis pipeline: " + err.Error(),
			Type:    "error",
			Entity:  "signal",
			Data:    fmt.Sprintf("Symbol: %s, Timeframe: %s", symbol, timeframe),
		})
		return res, fmt.Errorf("Redis read error: %v", err)
	}

	var lastSignal string
	lastSignalRes := rdb.Get(context.Background(), lastSignalKey)
	if lastSignalRes.Err() != nil {
		lastSignal = consts.HoldSignal
	} else {
		lastSignal = lastSignalRes.Val()
	}

	vals50, _ := res50.Result()
	vals200, _ := res200.Result()

	if len(vals50) < 50 || len(vals200) < 200 {
		return res, nil
	}

	ma50 := average(vals50)
	ma200 := average(vals200)

	res.MA50 = ma50.String()
	res.MA200 = ma200.String()

	if lastSignal == consts.BuySignal && ma50.GreaterThan(ma200) {
		return res, errors.New(consts.AlreadyInBuy)
	} else if lastSignal == consts.SellSignal && ma50.LessThan(ma200) {
		return res, errors.New(consts.AlreadyInSell)
	}

	jsonData, _ := json.Marshal(
		map[string]interface{}{
			"MA50":  res.MA50,
			"MA200": res.MA200,
		},
	)

	signal := entities.Signal{
		Symbol:    symbol,
		Timeframe: timeframe,
		Indicator: string(jsonData),
		LastTrade: "{}",
	}

	if ma50.GreaterThan(ma200) {
		rdb.Set(context.Background(), lastSignalKey, consts.BuySignal, 0)
		res.Signal = consts.BuySignal
		signal.Signal = consts.BuySignal
	}

	if ma50.LessThan(ma200) {
		rdb.Set(context.Background(), lastSignalKey, consts.SellSignal, 0)
		res.Signal = consts.SellSignal
		signal.Signal = consts.SellSignal
	}

	db := database.PgClient()
	if err := db.Create(&signal).Error; err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting signal into Postgres",
			Message: "Error inserting signal into Postgres: " + err.Error(),
			Type:    "error",
			Entity:  "signal",
			Data:    fmt.Sprintf("Symbol: %s, Signal: %s", signal.Symbol, signal.Signal),
		})
		log.Printf("Error inserting signal into Postgres: %v", err)
	}

	//add to mongo signal data
	mongoDb := database.MongoClient()
	collName := consts.CollectionNameSignal
	collection := mongoDb.Database(config.ReadValue().Mongo.Database).Collection(collName)
	_, err = collection.InsertOne(context.Background(), signal)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error inserting signal into MongoDB",
			Message: "Error inserting signal into MongoDB: " + err.Error(),
			Type:    "error",
			Entity:  "signal",
			Data:    fmt.Sprintf("Symbol: %s, Signal: %s", signal.Symbol, signal.Signal),
		})
		log.Printf("Error inserting signal into MongoDB: %v", err)
	}

	return res, errors.New(consts.AlreadyInHold)
}

func average(values []string) decimal.Decimal {
	var sum decimal.Decimal
	for _, val := range values {
		f, err := decimal.NewFromString(val)
		if err != nil {
			continue
		}
		sum = sum.Add(f)
	}
	return sum.Div(decimal.NewFromInt(int64(len(values))))
}
