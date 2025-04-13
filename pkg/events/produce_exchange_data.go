package events

import (
	"fmt"
	"log"
	"strings"

	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Stream struct {
	DB    *gorm.DB
	Kafka *KafkaClient
}

func NewStream(db *gorm.DB, kafkaClient *KafkaClient) *Stream {
	return &Stream{
		DB:    db,
		Kafka: kafkaClient,
	}
}

func (s *Stream) GetExchanges() []entities.Exchange {
	var exchanges []entities.Exchange
	err := s.DB.Where("is_active = ?", entities.ExchangeActive).Find(&exchanges).Error
	if err != nil {
		return exchanges
	}
	return exchanges
}

func (s *Stream) GetStreamWS(exchangeID string) string {
	var exchange entities.Exchange
	err := s.DB.Where("id = ?", exchangeID).First(&exchange).Error
	if err != nil {
		return ""
	}
	return exchange.WsUrl
}

func (s *Stream) GetStreamSymbols(exchangeID string) ([]entities.Symbol, error) {
	var symbols []entities.Symbol
	err := s.DB.Where("exchange_id = ?", exchangeID).Find(&symbols).Error
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

func (s *Stream) GetSymbolIntervals(exchangeID, symbol string) ([]entities.SignalInterval, error) {

	var intervals []entities.SignalInterval
	err := s.DB.Debug().Where("exchange_id = ? AND symbol = ?", exchangeID, strings.ToLower(symbol)).Find(&intervals).Error
	if err != nil {
		return intervals, err
	}
	return intervals, nil
}

func (s *Stream) StartAllStreams(exchangeID, topic string) error {
	wsBase := s.GetStreamWS(exchangeID)
	if wsBase == "" {
		return fmt.Errorf("WebSocket URL not found for exchange ID: %s", exchangeID)
	}

	symbols, err := s.GetStreamSymbols(exchangeID)
	if err != nil {
		return err
	}

	for _, symbol := range symbols {

		switch topic {
		case consts.OrderBookTopic:

			go func() {

				wsURL := fmt.Sprintf("%s/ws/%s@%s", wsBase, strings.ToLower(symbol.Symbol), consts.StreamOrderBook)
				log.Printf("Connecting to WS for %s symbol: %s", consts.StreamOrderBook, symbol.Symbol)

				err := s.startSymbolStream(wsURL, symbol.Symbol, topic)
				if err != nil {
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
				log.Printf("Error getting intervals for %s: %v", symbol.Symbol, err)
				return err
			}

			if len(intervals) == 0 {
				log.Printf("No intervals found for %s", symbol.Symbol)
				return err
			}

			for _, interval := range intervals {
				go func() {
					wsURL := fmt.Sprintf("%s/ws/%s@%s", wsBase, strings.ToLower(symbol.Symbol), fmt.Sprintf(consts.StreamCandleStick, interval.Interval))
					log.Printf("Connecting to WS for interval: %s symbol: %s", fmt.Sprintf(consts.StreamCandleStick, interval.Interval), symbol.Symbol)
					err := s.startSymbolStream(wsURL, symbol.Symbol, topic)
					if err != nil {
						log.Printf("Error in stream for %s: %v", symbol.Symbol, err)
					}

				}()
			}

		default:
			log.Printf("Unknown topic: %s", topic)
			continue
		}

	}

	return nil
}

func (s *Stream) startSymbolStream(wsURL, symbol, topic string) error {
	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("WebSocket dial failed for %s: %v", symbol, err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("[%s] Read error: %v", symbol, err)
			break
		}
		//log.Printf("[%s] Received message for %s", symbol, topic)

		//log.Println("topic", topic)
		_, _, err = s.Kafka.Produce(topic, symbol, []byte(message))
		if err != nil {
			log.Printf("[%s] Kafka write error: %v", symbol, err)
		}
	}

	return nil
}
