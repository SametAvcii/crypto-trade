package exchange

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

func (m *MockRepository) AddExchange(ctx context.Context, req dtos.AddExchangeReq) (dtos.AddExchangeRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.AddExchangeRes), args.Error(1)
}

func (m *MockRepository) UpdateExchange(ctx context.Context, req dtos.UpdateExchangeReq) (dtos.UpdateExchangeRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.UpdateExchangeRes), args.Error(1)
}

func (m *MockRepository) GetExchangeById(ctx context.Context, id string) (dtos.GetExchangeRes, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.GetExchangeRes), args.Error(1)
}

func (m *MockRepository) DeleteExchange(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetAllExchanges(ctx context.Context) ([]dtos.GetExchangeRes, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dtos.GetExchangeRes), args.Error(1)
}

func TestAddExchange(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := dtos.AddExchangeReq{Name: "Test Exchange"}
	expected := dtos.AddExchangeRes{ID: "1", Name: "Test Exchange"}

	mockRepo.On("AddExchange", mock.Anything, req).Return(expected, nil)

	result, err := service.AddExchange(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := dtos.UpdateExchangeReq{ID: "1", Name: "Updated Exchange"}
	expected := dtos.UpdateExchangeRes{ID: "1", Name: "Updated Exchange"}

	mockRepo.On("UpdateExchange", mock.Anything, req).Return(expected, nil)

	result, err := service.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetById(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := dtos.GetExchangeRes{ID: "1", Name: "Test Exchange"}

	mockRepo.On("GetExchangeById", mock.Anything, "1").Return(expected, nil)

	result, err := service.GetById(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("DeleteExchange", mock.Anything, "1").Return(nil)

	err := service.Delete(context.Background(), "1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := []dtos.GetExchangeRes{
		{ID: "1", Name: "Exchange 1"},
		{ID: "2", Name: "Exchange 2"},
	}

	mockRepo.On("GetAllExchanges", mock.Anything).Return(expected, nil)

	result, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
