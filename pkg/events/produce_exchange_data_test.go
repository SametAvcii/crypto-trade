package events_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SametAvcii/crypto-trade/internal/clients/kafka"
	"github.com/SametAvcii/crypto-trade/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	mockKafka := &kafka.KafkaClient{}

	stream := events.NewStream(mockDB, mockKafka)

	assert.NotNil(t, stream)
	assert.Equal(t, mockDB, stream.DB)
	assert.Equal(t, mockKafka, stream.Kafka)
}

func TestGetStreamSymbols(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	// Step 3: expected query & return rows
	mock.ExpectQuery(`SELECT \* FROM "symbols" WHERE exchange_id = \$1`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").
		WillReturnRows(sqlmock.NewRows([]string{"id", "symbol", "exchange_id"}).
			AddRow("750e8400-e29b-41d4-a716-446655440000", "BTCUSDT", "550e8400-e29b-41d4-a716-446655440000").
			AddRow("650e8400-e29b-41d4-a716-446655440000", "ETHUSDT", "550e8400-e29b-41d4-a716-446655440000"))

	// Step 4: call function
	stream := &events.Stream{
		DB:    gormDB,
		Kafka: &kafka.KafkaClient{},
	}

	symbols, err := stream.GetStreamSymbols("550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)
	assert.Len(t, symbols, 2)
	assert.Equal(t, "BTCUSDT", symbols[0].Symbol)
	assert.Equal(t, "ETHUSDT", symbols[1].Symbol)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSymbolIntervals(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Step 2: Wrap with GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	exchangeID := "550e8400-e29b-41d4-a716-446655440000"
	symbol := "BTCUSDT"

	// Step 3: Set up expected query and return rows
	query := regexp.QuoteMeta(`SELECT * FROM "signal_intervals" WHERE (exchange_id = $1 AND symbol = $2) AND "signal_intervals"."deleted_at" IS NULL`)
	mock.ExpectQuery(query).
		WithArgs(exchangeID, strings.ToLower(symbol)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "exchange_id", "symbol", "interval"}).
			AddRow("850e8400-e29b-41d4-a716-446655440000", exchangeID, "btcusdt", "1m").
			AddRow("950e8400-e29b-41d4-a716-446655440000", exchangeID, "btcusdt", "5m"))

	// Step 4: Create stream object and call function
	stream := &events.Stream{
		DB:    gormDB,
		Kafka: &kafka.KafkaClient{},
	}

	intervals, err := stream.GetSymbolIntervals(exchangeID, symbol)

	// Step 5: Assert expectations
	assert.NoError(t, err)
	assert.Len(t, intervals, 2)
	assert.Equal(t, "1m", intervals[0].Interval)
	assert.Equal(t, "5m", intervals[1].Interval)
	assert.Equal(t, "btcusdt", intervals[0].Symbol)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
