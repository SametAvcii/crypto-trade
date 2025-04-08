package events

import (
	"fmt"
	"log"
	"strings"

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

func (s *Stream) StartAllStreams(exchangeID string) error {
	log.Println("Starting streams for exchange:", exchangeID)
	wsBase := s.GetStreamWS(exchangeID)
	if wsBase == "" {
		return fmt.Errorf("WebSocket URL not found for exchange ID: %s", exchangeID)
	}

	symbols, err := s.GetStreamSymbols(exchangeID)
	if err != nil {
		return err
	}

	for _, sym := range symbols {
		symbol := sym

		go func() {

			
			wsURL := fmt.Sprintf("%s/ws/%s@aggTrade", wsBase, strings.ToLower(symbol.Symbol))
			log.Printf("Connecting to WS for symbol: %s", symbol.Symbol)

			err := startSymbolStream(wsURL, symbol.Symbol, s.Kafka)
			if err != nil {
				log.Printf("Error in stream for %s: %v", symbol.Symbol, err)
			}
		}()
	}

	return nil
}

func startSymbolStream(wsURL, symbol string, writer *KafkaClient) error {
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

		_, _, err = writer.Produce("crypto-trade", symbol, []byte(message))
		if err != nil {
			log.Printf("[%s] Kafka write error: %v", symbol, err)
		}
	}

	return nil
}
