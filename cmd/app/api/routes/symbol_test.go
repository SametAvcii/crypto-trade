package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSymbolService struct {
	mock.Mock
}

func (m *MockSymbolService) AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return dtos.AddSymbolRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.AddSymbolRes), args.Error(1)
}

func (m *MockSymbolService) UpdateSymbol(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return dtos.UpdateSymbolRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.UpdateSymbolRes), args.Error(1)
}

func (m *MockSymbolService) DeleteSymbol(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSymbolService) GetAllSymbols(ctx context.Context) ([]dtos.GetSymbolRes, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dtos.GetSymbolRes), args.Error(1)
}

func (m *MockSymbolService) GetSymbol(ctx context.Context, id string) (dtos.GetSymbolRes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return dtos.GetSymbolRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.GetSymbolRes), args.Error(1)
}

func TestAddSymbol(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	SymbolRoutes(group, mockService)

	t.Run("success", func(t *testing.T) {
		req := dtos.AddSymbolReq{Symbol: "BTC/USDT"}
		expectedRes := &dtos.GetSymbolRes{Symbol: "BTC/USDT"}

		mockService.On("AddSymbol", mock.Anything, req).Return(expectedRes, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/symbol/", bytes.NewBuffer(body))
		router.ServeHTTP(w, r)

		assert.Equal(t, 201, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		req := dtos.AddSymbolReq{Symbol: ""}
		mockService.On("AddSymbol", mock.Anything, req).Return(nil, errors.New("invalid symbol"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/symbol/", bytes.NewBuffer(body))
		router.ServeHTTP(w, r)

		assert.Equal(t, 400, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetSymbolByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	SymbolRoutes(group, mockService)

	t.Run("success", func(t *testing.T) {
		expectedRes := &dtos.GetSymbolRes{Symbol: "BTC/USDT"}
		mockService.On("GetSymbol", mock.Anything, "1").Return(expectedRes, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/symbol/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetSymbol", mock.Anything, "999").Return(nil, errors.New("symbol not found"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/symbol/999", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, 400, w.Code)
		mockService.AssertExpectations(t)
	})
}
