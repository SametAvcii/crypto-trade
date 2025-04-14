package symbol

import (
	"context"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddSymbol(ctx context.Context, req dtos.AddSymbolReq) (dtos.AddSymbolRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.AddSymbolRes), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (dtos.GetSymbolRes, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.GetSymbolRes), args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context) ([]dtos.GetSymbolRes, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dtos.GetSymbolRes), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, req dtos.UpdateSymbolReq) (dtos.UpdateSymbolRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.UpdateSymbolRes), args.Error(1)
}

func TestAddSymbol(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := dtos.AddSymbolReq{Symbol: "BTC/USDT"}
	expected := &dtos.AddSymbolRes{ID: "1", Symbol: "BTC/USDT"}

	mockRepo.On("AddSymbol", mock.Anything, req).Return(expected, nil)

	result, err := service.AddSymbol(t.Context(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSymbol(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := &dtos.GetSymbolRes{ID: "1", Symbol: "BTC/USDT"}

	mockRepo.On("GetByID", mock.Anything, "1").Return(expected, nil)

	result, err := service.GetSymbol(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllSymbols(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := []dtos.GetSymbolRes{
		{ID: "1", Symbol: "BTC/USDT"},
		{ID: "2", Symbol: "ETH/USDT"},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

	result, err := service.GetAllSymbols(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSymbol(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "1").Return(nil)

	err := service.DeleteSymbol(t.Context(), "1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSymbol(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := dtos.UpdateSymbolReq{ID: "1", Symbol: "Updated BTC/USDT"}
	expected := &dtos.UpdateSymbolRes{ID: "1", Symbol: "Updated BTC/USDT"}

	mockRepo.On("Update", mock.Anything, req).Return(expected, nil)

	result, err := service.UpdateSymbol(t.Context(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
