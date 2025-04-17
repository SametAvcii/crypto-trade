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
		req := dtos.AddSymbolReq{Symbol: "BTC/USDT", ExchangeID: "1"}
		expectedRes := dtos.AddSymbolRes{Symbol: "BTC/USDT", ExchangeID: "1"}

		mockService.On("AddSymbol", mock.Anything, req).Return(expectedRes, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/symbol/", bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equal(t, 201, w.Code)

		var response struct {
			Data   dtos.AddSymbolRes `json:"data"`
			Status int               `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, response.Data)
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

func TestUpdateSymbol(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	group.PUT("/:id", UpdateSymbol(mockService)) // Register handler

	t.Run("success", func(t *testing.T) {
		req := dtos.UpdateSymbolReq{
			Symbol:     "ETH/USDT",
			ExchangeID: "2",
		}
		id := "123"
		expectedRes := dtos.UpdateSymbolRes{
			ID:         id,
			Symbol:     req.Symbol,
			ExchangeID: req.ExchangeID,
		}

		// The handler sets req.ID = id from the path
		expectedReq := req
		expectedReq.ID = id

		mockService.On("UpdateSymbol", mock.Anything, expectedReq).Return(expectedRes, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/symbol/"+id, bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data   dtos.UpdateSymbolRes `json:"data"`
			Status int                  `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("bind error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/symbol/456", bytes.NewBufferString(`bad-json`))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		req := dtos.UpdateSymbolReq{
			Symbol:     "ETH/USDT",
			ExchangeID: "2",
		}
		id := "456"
		expectedReq := req
		expectedReq.ID = id

		mockService.On("UpdateSymbol", mock.Anything, expectedReq).Return(dtos.UpdateSymbolRes{}, errors.New("update failed"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/symbol/"+id, bytes.NewBuffer(body))
		r.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestDeleteSymbol(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	group.DELETE("/:id", DeleteSymbol(mockService)) // Register route

	t.Run("success", func(t *testing.T) {
		id := "123"

		mockService.On("DeleteSymbol", mock.Anything, id).Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/symbol/"+id, nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Message string `json:"message"`
			Status  int    `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Symbol deleted successfully", response.Message)
		assert.Equal(t, 200, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		id := "456"
		mockService.On("DeleteSymbol", mock.Anything, id).Return(errors.New("delete failed"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/symbol/"+id, nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response struct {
			Error  string `json:"error"`
			Status int    `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "delete failed", response.Error)
		assert.Equal(t, 400, response.Status)

		mockService.AssertExpectations(t)
	})
}
func TestGetAllSymbols(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	group.GET("", GetAllSymbols(mockService)) // Register route

	t.Run("success", func(t *testing.T) {
		expectedRes := []dtos.GetSymbolRes{
			{ID: "1", Symbol: "BTC/USDT", ExchangeID: "1"},
			{ID: "2", Symbol: "ETH/USDT", ExchangeID: "2"},
		}

		mockService.On("GetAllSymbols", mock.Anything).Return(expectedRes, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/symbol", nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data   []dtos.GetSymbolRes `json:"data"`
			Status int                 `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, response.Data)
		assert.Equal(t, 200, response.Status)

		mockService.AssertExpectations(t)
	})

}

func TestGetSymbolById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockSymbolService)
	router := gin.Default()
	group := router.Group("/symbol")
	group.GET("/:id", GetSymbolByID(mockService)) // Register route

	t.Run("success", func(t *testing.T) {
		id := "123"
		expectedRes := dtos.GetSymbolRes{ID: id, Symbol: "BTC/USDT", ExchangeID: "1"}

		mockService.On("GetSymbol", mock.Anything, id).Return(expectedRes, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/symbol/"+id, nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data   dtos.GetSymbolRes `json:"data"`
			Status int               `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedRes, response.Data)
		assert.Equal(t, 200, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		id := "456"
		mockService.On("GetSymbol", mock.Anything, id).Return(dtos.GetSymbolRes{}, errors.New("error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/symbol/"+id, nil)

		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response struct {
			Error  string `json:"error"`
			Status int    `json:"status"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Error)
		assert.Equal(t, 400, response.Status)

		mockService.AssertExpectations(t)
	})
}
