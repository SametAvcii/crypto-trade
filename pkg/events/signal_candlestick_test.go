package events

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDBSignal is a mock database client
type MockDBSignal struct {
	mock.Mock
}


func TestHandleMessage_UnclosedCandlestick(t *testing.T) {
	// Setup test data with unclosed candlestick
	kline := dtos.Kline{
		IsKlineClosed: false,
		ClosePrice:    "100.0",
		Interval:      "1h",
	}
	payload := dtos.CandlestickWs{
		Symbol: "BTCUSDT",
		Kline:  kline,
	}
	msgValue, _ := json.Marshal(payload)
	msg := &sarama.ConsumerMessage{
		Value: msgValue,
	}

	// Execute
	handler := SignalHandlerCandleStick{}
	handler.HandleMessage(msg)

	// No assertion needed as the function should just return early
}

func TestAverage(t *testing.T) {
	testCases := []struct {
		values   []string
		expected string
	}{
		{
			values:   []string{"10.0", "20.0", "30.0"},
			expected: "20",
		},
		{
			values:   []string{"1.5", "2.5", "3.5", "4.5"},
			expected: "3",
		},
		{
			values:   []string{},
			expected: "0", // Empty array should return 0
		},
		{
			values:   []string{"invalid", "20.0", "30.0"},
			expected: "25", // Invalid values should be skipped
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			result := average(tc.values)
			assert.Equal(t, tc.expected, result.String())
		})
	}
}

// Helper function to generate test values
func generateValues(count int, value string) []string {
	values := make([]string, count)
	for i := 0; i < count; i++ {
		values[i] = value
	}
	return values
}

func TestCheckForSignal(t *testing.T) {
	// Setup Redis mock
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Run("Test Buy Signal", func(t *testing.T) {
		// Setup test data for buy signal (MA50 > MA200)
		symbol := "BTCUSDT"
		timeframe := "1h"

		// Populate Redis with test data
		ma50Values := generateValues(50, "200.0")   // Higher values
		ma200Values := generateValues(200, "100.0") // Lower values

		key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
		key200 := fmt.Sprintf("%s:%s:ma200", symbol, timeframe)

		for _, v := range ma50Values {
			rdb.RPush(context.Background(), key50, v)
		}
		for _, v := range ma200Values {
			rdb.RPush(context.Background(), key200, v)
		}

		result, err := checkForSignal(rdb, symbol, timeframe)
		assert.NoError(t, err)
		assert.Equal(t, consts.BuySignal, result.Signal)
	})

	t.Run("Test Sell Signal", func(t *testing.T) {
		// Clear previous data
		rdb.FlushAll(context.Background())

		symbol := "BTCUSDT"
		timeframe := "1h"

		// Populate Redis with test data for sell signal (MA50 < MA200)
		ma50Values := generateValues(50, "100.0")   // Lower values
		ma200Values := generateValues(200, "200.0") // Higher values

		key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
		key200 := fmt.Sprintf("%s:%s:ma200", symbol, timeframe)

		for _, v := range ma50Values {
			rdb.RPush(context.Background(), key50, v)
		}
		for _, v := range ma200Values {
			rdb.RPush(context.Background(), key200, v)
		}

		result, err := checkForSignal(rdb, symbol, timeframe)
		assert.NoError(t, err)
		assert.Equal(t, consts.SellSignal, result.Signal)
	})

	t.Run("Test Insufficient Data", func(t *testing.T) {
		// Clear previous data
		rdb.FlushAll(context.Background())

		symbol := "BTCUSDT"
		timeframe := "1h"

		// Add insufficient data
		key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
		rdb.RPush(context.Background(), key50, "100.0")

		result, err := checkForSignal(rdb, symbol, timeframe)
		assert.NoError(t, err)
		assert.Empty(t, result.Signal)
	})

	t.Run("Test Already in Buy", func(t *testing.T) {
		// Clear previous data
		rdb.FlushAll(context.Background())

		symbol := "BTCUSDT"
		timeframe := "1h"

		// Set last signal as buy
		lastSignalKey := fmt.Sprintf("%s:%s:lastSignal", symbol, timeframe)
		rdb.Set(context.Background(), lastSignalKey, consts.BuySignal, 0)

		// Setup MA50 > MA200
		ma50Values := generateValues(50, "200.0")
		ma200Values := generateValues(200, "100.0")

		key50 := fmt.Sprintf("%s:%s:ma50", symbol, timeframe)
		key200 := fmt.Sprintf("%s:%s:ma200", symbol, timeframe)

		for _, v := range ma50Values {
			rdb.RPush(context.Background(), key50, v)
		}
		for _, v := range ma200Values {
			rdb.RPush(context.Background(), key200, v)
		}

		_, err := checkForSignal(rdb, symbol, timeframe)
		assert.EqualError(t, err, consts.AlreadyInBuy)
	})
}
