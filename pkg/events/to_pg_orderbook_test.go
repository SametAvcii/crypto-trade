package events

import (
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestPgOrderBookHandler_HandleMessage(t *testing.T) {
	// Setup test Redis and DB connections
	cache.InitRedis(config.ReadValue().Redis)
	database.InitDB(config.ReadValue().Database)

	handler := &PgOrderBookHandler{}

	tests := []struct {
		name    string
		message *sarama.ConsumerMessage
		wantErr bool
	}{
		{
			name: "Valid orderbook message",
			message: &sarama.ConsumerMessage{
				Value: []byte(`{
					"symbol": "BTCUSDT",
					"bids": [["50000.00", "1.00000000"]],
					"asks": [["51000.00", "0.50000000"]]
				}`),
			},
			wantErr: false,
		},
		{
			name: "Invalid JSON message",
			message: &sarama.ConsumerMessage{
				Value: []byte(`invalid json`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.HandleMessage(tt.message)
		})
	}
}

func TestUpdateOrderBookData(t *testing.T) {
	cache.InitRedis(config.ReadValue().Redis)
	database.InitDB(config.ReadValue().Database)

	rdb := cache.RedisClient()
	ctx := context.Background()

	tests := []struct {
		name    string
		symbol  string
		bids    [][]string
		asks    [][]string
		wantErr bool
	}{
		{
			name:   "Valid order book update",
			symbol: "BTCUSDT",
			bids:   [][]string{{"49000.00", "1.00000000"}},
			asks:   [][]string{{"51000.00", "0.50000000"}},
		},
		{
			name:   "Update with empty orders",
			symbol: "ETHUSDT",
			bids:   [][]string{},
			asks:   [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateOrderBookData(tt.symbol, tt.bids, tt.asks)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify Redis data
			bidKey := "order-book-depth:" + tt.symbol + ":bids"
			askKey := "order-book-depth:" + tt.symbol + ":asks"

			storedBids, _ := rdb.HGetAll(ctx, bidKey).Result()
			storedAsks, _ := rdb.HGetAll(ctx, askKey).Result()

			assert.Equal(t, len(tt.bids), len(storedBids))
			assert.Equal(t, len(tt.asks), len(storedAsks))
		})
	}
}

func TestCompareAndUpdate(t *testing.T) {
	database.InitDB(config.ReadValue().Database)
	tests := []struct {
		name     string
		oldData  map[string]string
		newData  [][]string
		symbol   string
		side     string
		expected int
	}{
		{
			name: "New orders added",
			oldData: map[string]string{
				"48000.00": "1.00000000",
			},
			newData: [][]string{
				{"48000.00", "1.00000000"},
				{"49000.00", "2.00000000"},
			},
			symbol:   "BTCUSDT",
			side:     "bid",
			expected: 2,
		},
		{
			name: "Order removed",
			oldData: map[string]string{
				"48000.00": "1.00000000",
			},
			newData: [][]string{
				{"48000.00", "0.00000000"},
			},
			symbol:   "BTCUSDT",
			side:     "bid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compareAndUpdate(tt.side, tt.oldData, tt.newData, tt.symbol)

		})
	}
}
