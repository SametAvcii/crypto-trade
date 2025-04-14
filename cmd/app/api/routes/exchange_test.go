package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExchangeService struct {
	mock.Mock
}

func (m *mockExchangeService) AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.AddExchangeRes), args.Error(1)
}

func (m *mockExchangeService) Update(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.UpdateExchangeRes), args.Error(1)
}

func (m *mockExchangeService) GetById(ctx context.Context, id string) (dtos.GetExchangeRes, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.GetExchangeRes), args.Error(1)
}

func (m *mockExchangeService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockExchangeService) GetAll(ctx context.Context) ([]dtos.GetExchangeRes, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dtos.GetExchangeRes), args.Error(1)
}

func TestAddExchange(t *testing.T) {
	mockService := new(mockExchangeService)
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	req := dtos.AddExchangeReq{
		Name: "Test Exchange",
	}
	res := dtos.AddExchangeRes{
		ID:   "1",
		Name: "Test Exchange",
	}

	mockService.On("AddExchange", mock.Anything, req).Return(res, nil)

	router.POST("/exchanges", AddExchange(mockService))

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/exchanges", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetExchangeById(t *testing.T) {
	mockService := new(mockExchangeService)
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	var exchange entities.Exchange

	newId := exchange.ID.String()
	exchange.Name = "Test Exchange"

	mockService.On("GetById", mock.Anything, newId).Return(exchange, nil)

	router.GET("/exchanges/:id", GetExchangeById(mockService))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/exchanges/1", nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteExchange(t *testing.T) {
	mockService := new(mockExchangeService)
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	mockService.On("Delete", mock.Anything, "1").Return(nil)

	router.DELETE("/exchanges/:id", DeleteExchange(mockService))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("DELETE", "/exchanges/1", nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllExchanges(t *testing.T) {
	mockService := new(mockExchangeService)
	router := gin.Default()
	gin.SetMode(gin.TestMode)

	exchanges := []dtos.GetExchangeRes{
		{ID: "1", Name: "Exchange 1"},
		{ID: "2", Name: "Exchange 2"},
	}

	mockService.On("GetAll", mock.Anything).Return(exchanges, nil)

	router.GET("/exchanges", GetAllExchanges(mockService))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/exchanges", nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
	mockService.AssertExpectations(t)
}
