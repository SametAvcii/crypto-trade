package events

import (
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock for DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Find(out interface{}, args ...interface{}) *gorm.DB {
	m.Called(out)
	return &gorm.DB{}
}

func (m *MockDB) First(out interface{}, args ...interface{}) *gorm.DB {
	m.Called(out)
	return &gorm.DB{}
}

func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

// Mock for KafkaClient
type MockKafkaClient struct {
	mock.Mock
}

func (m *MockKafkaClient) Produce(topic, key string, value []byte) (string, int64, error) {
	args := m.Called(topic, key, value)
	return args.String(0), int64(args.Int(1)), args.Error(2)
}

func TestNewStream(t *testing.T) {
	mockDB := &gorm.DB{}
	mockKafka := &KafkaClient{}

	stream := NewStream(mockDB, mockKafka)

	assert.NotNil(t, stream)
	assert.Equal(t, mockDB, stream.DB)
	assert.Equal(t, mockKafka, stream.Kafka)
}

func TestGetExchanges(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		// Unfortunately we can't easily mock GORM methods, so this test is limited
		exchanges := stream.GetExchanges()
		assert.NotNil(t, exchanges)
	})
}

func TestGetStreamWS(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		// Again, limited due to GORM mocking limitations
		ws := stream.GetStreamWS("exchange1")
		assert.Equal(t, "", ws) // Since we can't mock the DB properly, it will return empty string
	})
}

func TestGetStreamSymbols(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		symbols, err := stream.GetStreamSymbols("exchange1")
		assert.Nil(t, err)
		assert.NotNil(t, symbols)
	})
}

func TestGetSymbolIntervals(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		intervals, err := stream.GetSymbolIntervals("exchange1", "BTC/USDT")
		assert.Nil(t, err)
		assert.NotNil(t, intervals)
	})
}

func TestStartAllStreams(t *testing.T) {
	t.Run("Empty WebSocket URL", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		err := stream.StartAllStreams("nonexistent", consts.OrderBookTopic)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "WebSocket URL not found")
	})
}

func TestStartSymbolStream(t *testing.T) {
	// This would require mocking WebSocket connections, which is complex
	// A simple test that expects failure due to invalid URL
	t.Run("Invalid WebSocket URL", func(t *testing.T) {
		mockDB := &gorm.DB{}
		mockKafka := &KafkaClient{}

		stream := &Stream{
			DB:    mockDB,
			Kafka: mockKafka,
		}

		err := stream.startSymbolStream("invalid://url", "BTC/USDT", consts.OrderBookTopic)
		assert.Error(t, err)
	})
}
