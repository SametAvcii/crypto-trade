package events

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SametAvcii/crypto-trade/internal/clients/kafka"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Stream struct {
	DB    *gorm.DB
	Kafka *kafka.KafkaClient
}

func NewStream(db *gorm.DB, kafkaClient *kafka.KafkaClient) *Stream {
	return &Stream{
		DB:    db,
		Kafka: kafkaClient,
	}
}

func (s *Stream) GetExchanges() []entities.Exchange {
	var exchanges []entities.Exchange
	err := s.DB.Where("is_active = ?", entities.ExchangeActive).Find(&exchanges).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching exchanges from Postgres",
			Message: fmt.Sprintf("Error fetching exchanges from Postgres: %v", err),
			Type:    "error",
			Entity:  "stream",
			Data:    "Error fetching exchanges from Postgres",
		})
		return exchanges
	}
	return exchanges
}

func (s *Stream) GetStreamWS(exchangeID string) string {
	var exchange entities.Exchange
	err := s.DB.Where("id = ?", exchangeID).First(&exchange).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching exchange from Postgres",
			Message: fmt.Sprintf("Error fetching exchange from Postgres for exchange ID %s: %v", exchangeID, err),
			Type:    "error",
			Entity:  "stream",
			Data:    fmt.Sprintf("Exchange ID: %s", exchangeID),
		})

		return ""
	}
	return exchange.WsUrl
}

func (s *Stream) GetStreamSymbols(exchangeID string) ([]entities.Symbol, error) {
	var symbols []entities.Symbol
	err := s.DB.Where("exchange_id = ?", exchangeID).Find(&symbols).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching symbols from Postgres",
			Message: fmt.Sprintf("Error fetching symbols from Postgres for exchange ID %s: %v", exchangeID, err),
			Type:    "error",
			Entity:  "stream",
			Data:    fmt.Sprintf("Exchange ID: %s", exchangeID),
		})
		return nil, err
	}
	return symbols, nil
}

func (s *Stream) GetSymbolIntervals(exchangeID, symbol string) ([]entities.SignalInterval, error) {

	var intervals []entities.SignalInterval
	err := s.DB.Where("exchange_id = ? AND symbol = ?", exchangeID, strings.ToLower(symbol)).Find(&intervals).Error
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Error fetching intervals from Postgres",
			Message: fmt.Sprintf("Error fetching intervals from Postgres for symbol %s: %v", symbol, err),
			Type:    "error",
			Entity:  "stream",
			Data:    fmt.Sprintf("Exchange ID: %s, Symbol: %s", exchangeID, symbol),
		})

		return intervals, err
	}
	if len(intervals) == 0 {
		return intervals, fmt.Errorf("no intervals found for symbol %s", symbol)
	}
	return intervals, nil
}

func (s *Stream) StartAllStreams(exchangeID, topic string) error {
	wsBase := s.GetStreamWS(exchangeID)
	if wsBase == "" {
		ctlog.CreateLog(&entities.Log{
			Title:   "WebSocket URL Not Found",
			Message: fmt.Sprintf("WebSocket URL not found for exchange ID %s", exchangeID),
			Type:    "error",
			Entity:  "stream",
			Data:    fmt.Sprintf("Exchange ID: %s", exchangeID),
		})
		return fmt.Errorf("WebSocket URL not found for exchange ID: %s", exchangeID)
	}

	symbols, err := s.GetStreamSymbols(exchangeID)
	if err != nil {
		ctlog.CreateLog(&entities.Log{
			Title:   "Get Stream Symbols Error",
			Message: fmt.Sprintf("Error getting stream symbols for exchange ID %s: %v", exchangeID, err),
			Type:    "error",
			Entity:  "stream",
			Data:    fmt.Sprintf("Exchange ID: %s", exchangeID),
		})
		return err
	}

	for _, sym := range symbols {
		symbol := sym

		switch topic {
		case consts.OrderBookTopic:
			go func() {
				wsURL := fmt.Sprintf("%s/ws/%s@%s", wsBase, strings.ToLower(symbol.Symbol), consts.StreamOrderBook)
				log.Printf("Connecting to WS for %s symbol: %s", consts.StreamOrderBook, symbol.Symbol)

				err := s.startSymbolStream(wsURL, symbol.Symbol, topic)
				if err != nil {
					ctlog.CreateLog(&entities.Log{
						Title:   "WebSocket Connection Error",
						Message: fmt.Sprintf("WebSocket connection error for %s: %v", symbol.Symbol, err),
						Type:    "error",
						Entity:  "stream",
						Data:    fmt.Sprintf("WebSocket URL: %s", wsURL),
					})
					log.Printf("Error in stream for %s: %v", symbol.Symbol, err)
				}
			}()

		case consts.AggTradeTopic:
			go func() {

				wsURL := fmt.Sprintf("%s/ws/%s@%s", wsBase, strings.ToLower(symbol.Symbol), consts.StreamAggTrade)
				log.Printf("Connecting to WS for %s  symbol: %s", consts.StreamAggTrade, symbol.Symbol)

				err := s.startSymbolStream(wsURL, symbol.Symbol, topic)
				if err != nil {
					log.Printf("Error in stream for %s: %v", symbol.Symbol, err)
				}
			}()

		case consts.CandleStickTopic:

			log.Printf("Connecting to WS for intervals symbol: %s", symbol.Symbol)

			intervals, err := s.GetSymbolIntervals(exchangeID, symbol.Symbol)
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Error fetching intervals from Postgres",
					Message: fmt.Sprintf("Error fetching intervals from Postgres for symbol %s: %v", symbol.Symbol, err),
					Type:    "error",
					Entity:  "stream",
					Data:    fmt.Sprintf("Symbol: %s", symbol.Symbol),
				})
				log.Printf("Error getting intervals for %s: %v", symbol.Symbol, err)
				continue
			}

			if len(intervals) == 0 {
				ctlog.CreateLog(&entities.Log{
					Title:   "No intervals found",
					Message: fmt.Sprintf("No intervals found for symbol %s", symbol.Symbol),
					Type:    "error",
					Entity:  "stream",
					Data:    fmt.Sprintf("Symbol: %s", symbol.Symbol),
				})
				log.Printf("No intervals found for %s", symbol.Symbol)
				continue
			}

			for _, interval := range intervals {
				go func() {
					wsURL := fmt.Sprintf("%s/ws/%s@%s", wsBase, strings.ToLower(symbol.Symbol), fmt.Sprintf(consts.StreamCandleStick, interval.Interval))
					log.Printf("Connecting to WS for interval: %s symbol: %s", fmt.Sprintf(consts.StreamCandleStick, interval.Interval), symbol.Symbol)
					err := s.startSymbolStream(wsURL, symbol.Symbol, topic)
					if err != nil {
						ctlog.CreateLog(&entities.Log{
							Title:   "WebSocket Connection Error",
							Message: fmt.Sprintf("WebSocket connection error for %s: %v", symbol.Symbol, err),
							Type:    "error",
							Entity:  "stream",
							Data:    fmt.Sprintf("WebSocket URL: %s", wsURL),
						})

						log.Printf("Error in stream for %s: %v", symbol.Symbol, err)
					}

				}()
			}

		default:
			log.Printf("Unknown topic: %s", topic)
			ctlog.CreateLog(&entities.Log{
				Title:   "Unknown Topic",
				Message: fmt.Sprintf("Unknown topic: %s", topic),
				Type:    "error",
				Entity:  "stream",
				Data:    fmt.Sprintf("Topic: %s", topic),
			})
			continue
		}
	}

	return nil
}

func (s *Stream) startSymbolStream(wsURL, symbol, topic string) error {
	maxRetries := consts.MaxRetries
	retryDelay := consts.RetryDelay * time.Second

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			lastErr = fmt.Errorf("WebSocket dial failed for %s (attempt %d/%d): %v", symbol, attempt+1, maxRetries, err)
			ctlog.CreateLog(&entities.Log{
				Title:   "WebSocket Connection Error",
				Message: lastErr.Error(),
				Type:    "error",
				Entity:  "stream",
				Data:    fmt.Sprintf("WebSocket URL: %s", wsURL),
			})

			log.Println(lastErr.Error())
			time.Sleep(retryDelay)
			continue
		}

		// Successfully connected, start reading messages
		log.Printf("[%s] WebSocket connected, starting to read messages", symbol)
		defer c.Close()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "WebSocket Read Error",
					Message: fmt.Sprintf("WebSocket read error for %s: %v", symbol, err),
					Type:    "error",
					Entity:  "stream",
					Data:    fmt.Sprintf("WebSocket URL: %s", wsURL),
				})

				log.Printf("[%s] WebSocket read error: %v", symbol, err)
				break
			}

			_, _, err = s.Kafka.Produce(topic, symbol, []byte(message))
			if err != nil {
				ctlog.CreateLog(&entities.Log{
					Title:   "Kafka Write Error",
					Message: fmt.Sprintf("Kafka write error for %s: %v", symbol, err),
					Type:    "error",
					Entity:  "stream",
					Data:    fmt.Sprintf("Topic: %s, Symbol: %s", topic, symbol),
				})
				log.Printf("[%s] Kafka write error: %v", symbol, err)
			}

			log.Printf("[%s] Message sent to Kafka for %s", symbol, topic)
		}

		return nil
	}

	return lastErr
}
