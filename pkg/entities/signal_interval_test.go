package entities

import (
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignalInterval_FromDto(t *testing.T) {
	tests := []struct {
		name    string
		dto     *dtos.AddSignalIntervalReq
		want    *SignalInterval
		wantErr bool
	}{
		{
			name: "valid conversion",
			dto: &dtos.AddSignalIntervalReq{
				Symbol:     "BTCUSDT",
				Interval:   "1h",
				ExchangeId: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: &SignalInterval{
				Symbol:     "btcusdt",
				Interval:   "1h",
				ExchangeID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SignalInterval{}
			err := s.FromDto(tt.dto)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Symbol, s.Symbol)
			assert.Equal(t, tt.want.Interval, s.Interval)
			assert.Equal(t, tt.want.ExchangeID, s.ExchangeID)
		})
	}
}

func TestSignalInterval_UpdateFromDto(t *testing.T) {
	tests := []struct {
		name    string
		initial *SignalInterval
		dto     dtos.UpdateSignalIntervalReq
		want    *SignalInterval
		wantErr bool
	}{
		{
			name: "update all fields",
			initial: &SignalInterval{
				Symbol:     "btcusdt",
				Interval:   "1h",
				ExchangeID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			dto: dtos.UpdateSignalIntervalReq{
				Symbol:     "ETHUSDT",
				Interval:   "4h",
				ExchangeId: "550e8400-e29b-41d4-a716-446655440000",
			},
			want: &SignalInterval{
				Symbol:     "ethusdt",
				Interval:   "4h",
				ExchangeID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			wantErr: false,
		},
		{
			name: "partial update",
			initial: &SignalInterval{
				Symbol:     "btcusdt",
				Interval:   "1h",
				ExchangeID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			dto: dtos.UpdateSignalIntervalReq{
				Symbol: "ETHUSDT",
			},
			want: &SignalInterval{
				Symbol:     "ethusdt",
				Interval:   "1h",
				ExchangeID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.initial.UpdateFromDto(tt.dto)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Symbol, tt.initial.Symbol)
			assert.Equal(t, tt.want.Interval, tt.initial.Interval)
			assert.Equal(t, tt.want.ExchangeID, tt.initial.ExchangeID)
		})
	}
}
