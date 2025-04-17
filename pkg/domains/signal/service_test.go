package signal

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

func (m *MockRepository) AddSignalIntervals(ctx context.Context, req dtos.AddSignalIntervalReq) (dtos.AddSignalIntervalRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.AddSignalIntervalRes), args.Error(1)
}

func (m *MockRepository) UpdateSignalInterval(ctx context.Context, req dtos.UpdateSignalIntervalReq) (dtos.UpdateSignalIntervalRes, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dtos.UpdateSignalIntervalRes), args.Error(1)
}

func (m *MockRepository) DeleteSignalInterval(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetSignalInterval(ctx context.Context, id string) (dtos.GetSignalIntervalRes, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dtos.GetSignalIntervalRes), args.Error(1)
}

func (m *MockRepository) GetAllSignalIntervals(ctx context.Context) ([]dtos.GetSignalIntervalRes, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dtos.GetSignalIntervalRes), args.Error(1)
}
func TestAddSignalIntervals_Service(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := dtos.AddSignalIntervalReq{
		Symbol:     "Test Interval",
		Interval:   "1m",
		ExchangeId: "550e8400-e29b-41d4-a716-446655440000",
	}
	expected := dtos.AddSignalIntervalRes{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Symbol:   "test interval",
		Interval: "1m",
	}

	mockRepo.
		On("AddSignalIntervals", mock.Anything, req).
		Return(expected, nil)

	got, err := svc.AddSignalIntervals(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSignalIntervals(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	req := dtos.UpdateSignalIntervalReq{ID: "1", Symbol: "Updated Interval"}
	expected := dtos.UpdateSignalIntervalRes{ID: "1", Symbol: "Updated Interval"}

	mockRepo.On("UpdateSignalInterval", mock.Anything, req).Return(expected, nil)

	result, err := service.UpdateSignalIntervals(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSignalIntervals(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("DeleteSignalInterval", mock.Anything, "1").Return(nil)

	err := service.DeleteSignalIntervals(context.Background(), "1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSignalIntervalById(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := dtos.GetSignalIntervalRes{ID: "1", Symbol: "Test Interval"}

	mockRepo.On("GetSignalInterval", mock.Anything, "1").Return(expected, nil)

	result, err := service.GetSignalIntervalById(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
