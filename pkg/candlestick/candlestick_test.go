package candlestick_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) Get(url string, headers map[string]string, response interface{}) error {
	args := m.Called(url, headers)
	if res, ok := response.(*[][]interface{}); ok {
		*res = args.Get(0).([][]interface{})
	}
	return args.Error(1)
}

type MockPgDB struct {
	mock.Mock
}

func (m *MockPgDB) Where(query string, args ...interface{}) *MockPgDB {
	m.Called(query, args)
	return m
}
func (m *MockPgDB) First(dest interface{}) *MockPgDB {
	m.Called(dest)
	return m
}
func (m *MockPgDB) Create(value interface{}) *MockPgDB {
	m.Called(value)
	return m
}
func (m *MockPgDB) Error() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockPgDB) Find(out interface{}) *MockPgDB {
	m.Called(out)
	return m
}
func (m *MockPgDB) Limit(limit int) *MockPgDB {
	m.Called(limit)
	return m
}
func (m *MockPgDB) Order(value interface{}) *MockPgDB {
	m.Called(value)
	return m
}

// --- Test Case ---

func TestGetCandleSticksAndUpdate_Success(t *testing.T) {
	mockAPI := new(MockAPI)
	mockPg := new(MockPgDB)

	symbol := "BTCUSDT"
	interval := "1m"
	limit := 1
	exchangeId := "1"

	// Fake candlestick
	candle := [][]interface{}{
		{
			json.Number("1712837400000"),               // Open time
			"69000.0", "69100.0", "68900.0", "69050.0", // O H L C
			"100.0", json.Number("1712837459999"), "200", "50.0", "300.0", "0.0",
		},
	}

	mockAPI.On("Get", mock.Anything, mock.Anything).Return(candle, nil)

	// Set expectations for mocked DB
	mockPg.On("Where", "id= ?", []interface{}{exchangeId}).Return(mockPg)
	mockPg.On("First", mock.Anything).Return(mockPg)
	mockPg.On("Create", mock.Anything).Return(mockPg)
	mockPg.On("Error").Return(nil)
	mockPg.On("Where", "symbol = ? AND interval = ?", symbol, interval).Return(mockPg)
	mockPg.On("Limit", limit).Return(mockPg)
	mockPg.On("Order", "open_time desc").Return(mockPg)
	mockPg.On("Find", mock.Anything).Return(mockPg)

	// TODO: inject mocks into GetCandleSticksAndUpdate using dependency injection or interfaces

	// Example only: This call won't work without actual refactor for injection.
	// result, err := candlestick.GetCandleSticksAndUpdate(exchangeId, symbol, interval, limit)

	// assert.NoError(t, err)
	// assert.NotEmpty(t, result)

	// mockAPI.AssertExpectations(t)
	// mockPg.AssertExpectations(t)
}


