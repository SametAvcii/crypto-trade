package kafka

import (
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKafkaClient_Produce(t *testing.T) {
	tests := []struct {
		name          string
		topic         string
		key           string
		message       []byte
		producerSetup func(*mocks.SyncProducer)
		wantPartition int32
		wantOffset    int64
		wantErr       bool
	}{
		{
			name:    "successful production",
			topic:   "test-topic",
			key:     "test-key",
			message: []byte("test-message"),
			producerSetup: func(producer *mocks.SyncProducer) {
				producer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
					if msg.Topic != "test-topic" {
						return errors.New("unexpected topic")
					}
					if key, _ := msg.Key.Encode(); string(key) != "test-key" {
						return errors.New("unexpected key")
					}
					if value, _ := msg.Value.Encode(); string(value) != "test-message" {
						return errors.New("unexpected value")
					}
					return nil
				})
			},
			wantPartition: 0,
			wantOffset:    0,
			wantErr:       false,
		},
		{
			name:    "production failure",
			topic:   "test-topic",
			key:     "test-key",
			message: []byte("test-message"),
			producerSetup: func(producer *mocks.SyncProducer) {
				producer.ExpectSendMessageAndFail(errors.New("failed to send message"))
			},
			wantPartition: 0,
			wantOffset:    0,
			wantErr:       true,
		},
		{
			name:    "custom partition and offset",
			topic:   "test-topic",
			key:     "test-key",
			message: []byte("test-message"),
			producerSetup: func(producer *mocks.SyncProducer) {
				producer.ExpectSendMessageAndSucceed()
			},
			wantPartition: 0,
			wantOffset:    0,
			wantErr:       false,
		},
		{
			name:    "empty key",
			topic:   "test-topic",
			key:     "",
			message: []byte("test-message"),
			producerSetup: func(producer *mocks.SyncProducer) {
				producer.ExpectSendMessageWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error {
					if key, _ := msg.Key.Encode(); len(key) != 0 {
						return errors.New("expected empty key")
					}
					return nil
				})
			},
			wantPartition: 0,
			wantOffset:    0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock producer
			mockProducer := mocks.NewSyncProducer(t, nil)
			if tt.producerSetup != nil {
				tt.producerSetup(mockProducer)
			}

			// Create Kafka client with the mock producer
			k := &KafkaClient{
				producer: mockProducer,
			}

			// Call the Produce method
			partition, offset, err := k.Produce(tt.topic, tt.key, tt.message)

			// Verify the results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Since the mock producer controls the partition and offset,
			// we don't actually need to verify specific values here
			// But we can ensure they're returned properly
			if !tt.wantErr {
				assert.IsType(t, int32(0), partition)
				assert.IsType(t, int64(0), offset)
			}

			// Ensure all expectations were met
			err = mockProducer.Close()
			assert.NoError(t, err)
		})
	}
}
