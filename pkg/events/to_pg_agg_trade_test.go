package events

import (
	"encoding/json"
	"testing"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
)

func TestPgAggTradeHandler_HandleMessage(t *testing.T) {
	tests := []struct {
		name    string
		message *sarama.ConsumerMessage
		wantErr bool
	}{
		{
			name: "valid message",
			message: &sarama.ConsumerMessage{
				Value: mustMarshal(dtos.AggTradeMongo{
					MongoID: "123",
					Value: string(mustMarshal(dtos.AggTrade{
						Symbol: "BTCUSDT",
						Price:  "50000",
					})),
				}),
			},
			wantErr: false,
		},
		{
			name: "invalid json message",
			message: &sarama.ConsumerMessage{
				Value: []byte("invalid json"),
			},
			wantErr: true,
		},
		{
			name: "empty message",
			message: &sarama.ConsumerMessage{
				Value: nil,
			},
			wantErr: true,
		},
	}

	handler := &PgAggTradeHandler{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.HandleMessage(tt.message)
		})
	}
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
