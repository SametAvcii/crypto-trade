package entities

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"testing"
)

func TestSignal_FromDto(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.Signal
		want Signal
	}{
		{
			name: "Should convert DTO to entity successfully",
			dto: dtos.Signal{
				Symbol:        "BTC/USDT",
				Timeframe:    "1h",
				Signal:       "buy",
				IndicatorData: `{"rsi": 70, "macd": "bullish"}`,
				LastTrade:    `{"price": 50000, "side": "buy"}`,
			},
			want: Signal{
				Symbol:     "BTC/USDT",
				Timeframe: "1h",
				Signal:    "buy",
				Indicator: `{"rsi": 70, "macd": "bullish"}`,
				LastTrade: `{"price": 50000, "side": "buy"}`,
			},
		},
		{
			name: "Should handle empty values",
			dto:  dtos.Signal{},
			want: Signal{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Signal{}
			s.FromDto(tt.dto)

			if s.Symbol != tt.want.Symbol {
				t.Errorf("Symbol = %v, want %v", s.Symbol, tt.want.Symbol)
			}
			if s.Timeframe != tt.want.Timeframe {
				t.Errorf("Timeframe = %v, want %v", s.Timeframe, tt.want.Timeframe)
			}
			if s.Signal != tt.want.Signal {
				t.Errorf("Signal = %v, want %v", s.Signal, tt.want.Signal)
			}
			if s.Indicator != tt.want.Indicator {
				t.Errorf("Indicator = %v, want %v", s.Indicator, tt.want.Indicator)
			}
			if s.LastTrade != tt.want.LastTrade {
				t.Errorf("LastTrade = %v, want %v", s.LastTrade, tt.want.LastTrade)
			}
		})
	}
}