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

type MockSignalService struct {
	mock.Mock
}

func (m *MockSignalService) AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return dtos.AddSignalIntervalRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.AddSignalIntervalRes), args.Error(1)
}

func (m *MockSignalService) UpdateSignalIntervals(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return dtos.UpdateSignalIntervalRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.UpdateSignalIntervalRes), args.Error(1)
}

func (m *MockSignalService) DeleteSignalIntervals(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSignalService) GetSignalIntervalById(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return dtos.GetSignalIntervalRes{}, args.Error(1)
	}
	return args.Get(0).(dtos.GetSignalIntervalRes), args.Error(1)
}

func (m *MockSignalService) GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dtos.GetSignalIntervalRes), args.Error(1)
}

func TestAddSignalInterval(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSignalService)
	router := gin.Default()
	group := router.Group("/signal")
	SignalRoutes(group, mockService)

	t.Run("Success", func(t *testing.T) {
		req := dtos.AddSignalIntervalReq{Symbol: "BTC"}
		expectedRes := dtos.AddSignalIntervalRes{Symbol: "BTC"}

		mockService.On("AddSignalIntervals", mock.Anything, req).Return(expectedRes, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/signal/interval", bytes.NewBuffer(body))
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Error", func(t *testing.T) {
		req := dtos.AddSignalIntervalReq{Symbol: "BTC"}
		mockService.On("AddSignalIntervals", mock.Anything, req).Return(nil, errors.New("error")).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/signal/interval", bytes.NewBuffer(body))
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateSignalInterval(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSignalService)
	router := gin.Default()
	group := router.Group("/signal")
	SignalRoutes(group, mockService)
	t.Run("Success", func(t *testing.T) {
		req := dtos.UpdateSignalIntervalReq{ID: "1", Symbol: "BTC"}
		expectedRes := dtos.UpdateSignalIntervalRes{ID: "1", Symbol: "BTC"}

		mockService.On("UpdateSignalIntervals", mock.Anything, req).Return(expectedRes, nil).Once()

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/signal/interval/1", bytes.NewBuffer(body))
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestDeleteSignalInterval(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSignalService)
	router := gin.Default()
	group := router.Group("/signal")
	SignalRoutes(group, mockService)

	t.Run("Success", func(t *testing.T) {
		mockService.On("DeleteSignalIntervals", mock.Anything, "1").Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/signal/interval/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetSignalIntervalByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSignalService)
	router := gin.Default()
	group := router.Group("/signal")
	SignalRoutes(group, mockService)

	t.Run("Success", func(t *testing.T) {
		expectedRes := dtos.GetSignalIntervalRes{ID: "1", Symbol: "BTC"}
		mockService.On("GetSignalIntervalById", mock.Anything, "1").Return(expectedRes, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/signal/interval/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAllSignalIntervals(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockSignalService)
	router := gin.Default()
	group := router.Group("/signal")
	SignalRoutes(group, mockService)

	t.Run("Success", func(t *testing.T) {
		expectedRes := []dtos.GetSignalIntervalRes{{Symbol: "BTC"}}
		mockService.On("GetAllSignalIntervals", mock.Anything).Return(expectedRes, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/signal/interval", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
