package events

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageHandler struct {
	mock.Mock
}

func (m *MockMessageHandler) HandleMessage(msg *sarama.ConsumerMessage) {
	m.Called(msg)
}

func TestConsumer_Start(t *testing.T) {
	tests := []struct {
		name    string
		brokers []string
		groupID string
		topic   string
		wantErr bool
	}{
		{
			name:    "valid configuration",
			brokers: []string{"localhost:9092"},
			groupID: "test-group",
			topic:   "test-topic",
			wantErr: false,
		},
		{
			name:    "invalid broker",
			brokers: []string{"invalid:9092"},
			groupID: "test-group",
			topic:   "test-topic",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := &MockMessageHandler{}
			consumer := &Consumer{
				Brokers: tt.brokers,
				GroupID: tt.groupID,
				Topic:   tt.topic,
				Handler: mockHandler,
			}

			err := consumer.Start()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsumerGroupHandler_ConsumeClaim(t *testing.T) {
	mockHandler := &MockMessageHandler{}
	handler := &consumerGroupHandler{
		handler: mockHandler,
	}

	mockSession := &MockConsumerGroupSession{}
	mockClaim := &MockConsumerGroupClaim{}

	testMsg := &sarama.ConsumerMessage{
		Topic: "test-topic",
		Value: []byte("test message"),
	}

	mockHandler.On("HandleMessage", testMsg).Return()
	mockSession.On("MarkMessage", testMsg, "").Return()

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(mockClaim.messagesChan)
	}()

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	mockHandler.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

// Mock implementations for Sarama interfaces
type MockConsumerGroupSession struct {
	mock.Mock
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	return nil
}

func (m *MockConsumerGroupSession) MemberID() string {
	return ""
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	return 0
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *MockConsumerGroupSession) Context() context.Context {
	return context.Background()
}

func (m *MockConsumerGroupSession) Commit() {
	m.Called()
}

type MockConsumerGroupClaim struct {
	mock.Mock
	messagesChan chan *sarama.ConsumerMessage
}

func (m *MockConsumerGroupClaim) Topic() string {
	return "test-topic"
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	return 0
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	if m.messagesChan == nil {
		m.messagesChan = make(chan *sarama.ConsumerMessage, 1)
		m.messagesChan <- &sarama.ConsumerMessage{
			Topic: "test-topic",
			Value: []byte("test message"),
		}
	}
	return m.messagesChan
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	return m.Called().Get(0).(int64)
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	return m.Called().Get(0).(int64)
}
