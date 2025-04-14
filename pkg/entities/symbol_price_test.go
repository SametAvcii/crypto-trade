package entities

import (
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"testing"
)

func TestSymbolPrice_FromDto(t *testing.T) {
	tests := []struct {
		name string
		dto  *dtos.AggTrade
		want *SymbolPrice
	}{
		{
			name: "should convert dto to entity successfully",
			dto: &dtos.AggTrade{
				Symbol:       "BTCUSDT",
				Price:       "50000.00",
				Quantity:    "0.5",
				TradeID:     123456,
				TradeTime:   1634567890,
				IsBuyerMaker: true,
				EventTime:   1634567890,
				EventType:   "aggTrade",
			},
			want: &SymbolPrice{
				Symbol:       "BTCUSDT",
				Price:       "50000.00", 
				Quantity:    "0.5",
				TradeID:     123456,
				TradeTime:   1634567890,
				IsBuyerMaker: true,
				EventTime:   1634567890,
				EventType:   "aggTrade",
			},
		},
		{
			name: "should handle empty values",
			dto: &dtos.AggTrade{},
			want: &SymbolPrice{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SymbolPrice{}
			s.FromDto(tt.dto)

			if s.Symbol != tt.want.Symbol {
				t.Errorf("Symbol = %v, want %v", s.Symbol, tt.want.Symbol)
			}
			if s.Price != tt.want.Price {
				t.Errorf("Price = %v, want %v", s.Price, tt.want.Price)
			}
			if s.Quantity != tt.want.Quantity {
				t.Errorf("Quantity = %v, want %v", s.Quantity, tt.want.Quantity)
			}
			if s.TradeID != tt.want.TradeID {
				t.Errorf("TradeID = %v, want %v", s.TradeID, tt.want.TradeID)
			}
			if s.TradeTime != tt.want.TradeTime {
				t.Errorf("TradeTime = %v, want %v", s.TradeTime, tt.want.TradeTime)
			}
			if s.IsBuyerMaker != tt.want.IsBuyerMaker {
				t.Errorf("IsBuyerMaker = %v, want %v", s.IsBuyerMaker, tt.want.IsBuyerMaker)
			}
			if s.EventTime != tt.want.EventTime {
				t.Errorf("EventTime = %v, want %v", s.EventTime, tt.want.EventTime)
			}
			if s.EventType != tt.want.EventType {
				t.Errorf("EventType = %v, want %v", s.EventType, tt.want.EventType)
			}
		})
	}
}