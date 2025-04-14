package events

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/stretchr/testify/assert"
)

func TestGetPgTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "AggTrade topic",
			topic:    consts.AggTradeTopic,
			expected: consts.PgAggTradeTopic,
		},
		{
			name:     "OrderBook topic", 
			topic:    consts.OrderBookTopic,
			expected: consts.PgOrderBookTopic,
		},
		{
			name:     "Invalid topic",
			topic:    "invalid_topic",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPgTopic(tt.topic)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCollectionName(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected string
	}{
		{
			name:     "AggTrade collection",
			topic:    consts.AggTradeTopic,
			expected: consts.CollectionNameTrade,
		},
		{
			name:     "OrderBook collection",
			topic:    consts.OrderBookTopic,
			expected: consts.CollectionNameOrder,
		},
		{
			name:     "Invalid topic",
			topic:    "invalid_topic",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getCollectionName(tt.topic)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMongoHandlerHandleMessage(t *testing.T) {
	handler := &MongoHandler{}
	
	testCases := []struct {
		name    string
		message *sarama.ConsumerMessage
	}{
		{
			name: "Invalid JSON message",
			message: &sarama.ConsumerMessage{
				Topic: consts.AggTradeTopic,
				Value: []byte("invalid json"),
			},
		},
		{
			name: "Empty message",
			message: &sarama.ConsumerMessage{
				Topic: consts.AggTradeTopic,
				Value: []byte(""),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler.HandleMessage(tc.message)
		})
	}
}