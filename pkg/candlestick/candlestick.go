package candlestick

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SametAvcii/crypto-trade/internal/clients/database"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/SametAvcii/crypto-trade/pkg/utils"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCandleSticksAndUpdate(ctx context.Context, exchangeId, symbol string, interval string, limit int) ([]entities.Candlestick, error) {
	var (
		mongoClient = database.MongoClient()
		collection  = mongoClient.Database(config.ReadValue().Mongo.Database).Collection(consts.CollectionNameCandleStick)
		database    = database.PgClient()
		exchange    entities.Exchange
	)

	err := database.Where("id= ?", exchangeId).First(&exchange).Error
	if err != nil {
		log.Printf("Error fetching exchange from PostgreSQL: %v", err)
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching exchange from PostgreSQL",
			Message: "Error fetching exchange from PostgreSQL: " + err.Error(),
			Type:    "error",
			Entity:  "exchange",
			Data:    fmt.Sprintf("Exchange ID: %s, Error: %s", exchangeId, err.Error()),
		})
		return nil, err
	}

	api := utils.NewAPI(exchange.RestUrl)

	var url string
	switch exchange.Name {
	case consts.Binance:
		url = fmt.Sprintf("/klines?symbol=%s"+"&interval=%s&limit=%d", symbol, interval, limit)
	default:
		log.Printf("Exchange %s not supported", exchange.Name)
		return nil, fmt.Errorf("exchange %s not supported", exchange.Name)
	}
	var klines [][]interface{}

	err = api.Get(url, nil, &klines)
	if err != nil {
		log.Printf("Error getting candlestick data: %v", err)
		ctlog.CreateLog(&entities.Log{
			Title:   "Error getting candlestick data",
			Message: "Error getting candlestick data: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    url,
		})

		return nil, err
	}

	for _, k := range klines {
		openTime, _ := k[0].(json.Number).Int64()
		open, _ := decimal.NewFromString(k[1].(string))
		high, _ := decimal.NewFromString(k[2].(string))
		low, _ := decimal.NewFromString(k[3].(string))
		close, _ := decimal.NewFromString(k[4].(string))
		volume, _ := decimal.NewFromString(k[5].(string))
		quoteVolume, _ := decimal.NewFromString(k[7].(string))
		takerBuyBaseVolume, _ := decimal.NewFromString(k[9].(string))
		takerBuyQuoteVolume, _ := decimal.NewFromString(k[10].(string))
		ignore, _ := decimal.NewFromString(k[11].(string))
		closeTime, _ := k[6].(json.Number).Int64()
		numberOfTrades, _ := k[8].(json.Number).Int64()

		kline := dtos.CandlestickRest{
			Symbol:              symbol,
			Interval:            interval,
			ExchangeId:          exchange.ID.String(),
			OpenTime:            openTime,
			Open:                open,
			High:                high,
			Low:                 low,
			Close:               close,
			Volume:              volume,
			CloseTime:           closeTime,
			QuoteVolume:         quoteVolume,
			NumberOfTrades:      numberOfTrades,
			TakerBuyBaseVolume:  takerBuyBaseVolume,
			TakerBuyQuoteVolume: takerBuyQuoteVolume,
			Ignore:              ignore,
		}

		var lastCandlestick dtos.CandlestickRest

		filter := bson.M{"symbol": symbol, "interval": interval}
		opts := options.FindOne().SetSort(bson.D{{Key: "openTime", Value: -1}})
		err := collection.FindOne(ctx, filter, opts).Decode(&lastCandlestick)

		if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("Error fetching last candlestick: %v", err)
			ctlog.CreateLog(&entities.Log{
				Title:   "Error fetching last candlestick",
				Message: "Error fetching last candlestick: " + err.Error(),
				Type:    "error",
				Entity:  "candlestick",
				Data:    fmt.Sprintf("Symbol: %s, Error: %s", symbol, err.Error()),
			})
			return nil, err
		}

		if lastCandlestick.OpenTime < kline.OpenTime {
			_, err := collection.InsertOne(ctx, kline)
			if err != nil {
				log.Printf("Error inserting candlestick into MongoDB: %v", err)
				ctlog.CreateLog(&entities.Log{
					Title:   "Error inserting candlestick into MongoDB",
					Message: "Error inserting candlestick into MongoDB: " + err.Error(),
					Type:    "error",
					Entity:  "candlestick",
					Data:    fmt.Sprintf("Symbol: %s, Candlestick: %+v", symbol, kline),
				})
				return nil, err
			}
			var candlestick entities.Candlestick
			candlestick.FromDto(&kline)
			err = database.Create(&candlestick).Error
			if err != nil {
				log.Printf("Error inserting candlestick into PostgreSQL: %v", err)
				ctlog.CreateLog(&entities.Log{
					Title:   "Error inserting candlestick into PostgreSQL",
					Message: "Error inserting candlestick into PostgreSQL: " + err.Error(),
					Type:    "error",
					Entity:  "candlestick",
					Data:    fmt.Sprintf("Symbol: %s, Candlestick: %+v", symbol, kline),
				})
				return nil, err
			}
			log.Printf("Inserted candlestick into PostgreSQL: %+v", kline)
			ctlog.CreateLog(&entities.Log{
				Title:   "Inserted candlestick into PostgreSQL",
				Message: fmt.Sprintf("Inserted candlestick into PostgreSQL: %+v", kline),
				Type:    "info",
				Entity:  "candlestick",
				Data:    fmt.Sprintf("Symbol: %s, Candlestick: %+v", symbol, kline),
			})
		} else {
			log.Printf("Candlestick already exists in MongoDB: %+v", kline)
			ctlog.CreateLog(&entities.Log{
				Title:   "Candlestick already exists in MongoDB",
				Message: fmt.Sprintf("Candlestick already exists in MongoDB: %+v", kline),
				Type:    "info",
				Entity:  "candlestick",
				Data:    fmt.Sprintf("Symbol: %s, Candlestick: %+v", symbol, kline),
			})

		}
	}

	var candlesticks []entities.Candlestick
	err = database.Where("symbol = ? AND interval = ?", symbol, interval).Limit(limit).Order("open_time desc").Find(&candlesticks).Error
	if err != nil {
		log.Printf("Error fetching candlesticks from PostgreSQL: %v", err)
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching candlesticks from PostgreSQL",
			Message: "Error fetching candlesticks from PostgreSQL: " + err.Error(),
			Type:    "error",
			Entity:  "candlestick",
			Data:    fmt.Sprintf("Symbol: %s, Interval: %s", symbol, interval),
		})
		return nil, err
	}

	return candlesticks, nil
}
