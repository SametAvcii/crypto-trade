package entities

import (
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSymbol_FromDto(t *testing.T) {
	tests := []struct {
		name    string
		dto     *dtos.AddSymbolReq
		wantErr bool
	}{
		{
			name: "valid conversion",
			dto: &dtos.AddSymbolReq{
				Symbol:     "BTCUSDT",
				ExchangeID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "invalid uuid",
			dto: &dtos.AddSymbolReq{
				Symbol:     "BTCUSDT",
				ExchangeID: "invalid-uuid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Symbol{}
			err := s.FromDto(tt.dto)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.dto.Symbol, s.Symbol)
			}
		})
	}
}

func TestSymbol_ToDto(t *testing.T) {
	symbol := &Symbol{
		Symbol:     "ETHUSDT",
		ExchangeID: uuid.New(),
	}
	symbol.ID = uuid.New()

	dto := symbol.ToDto()
	assert.Equal(t, symbol.ID.String(), dto.ID)
	assert.Equal(t, symbol.Symbol, dto.Symbol)
	assert.Equal(t, symbol.ExchangeID.String(), dto.ExchangeID)
}

func TestSymbol_UpdateFromDto(t *testing.T) {
	tests := []struct {
		name    string
		dto     *dtos.UpdateSymbolReq
		wantErr bool
	}{
		{
			name: "valid update",
			dto: &dtos.UpdateSymbolReq{
				Symbol:     "BTCUSDT",
				ExchangeID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "invalid uuid",
			dto: &dtos.UpdateSymbolReq{
				Symbol:     "BTCUSDT",
				ExchangeID: "invalid-uuid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Symbol{}
			err := s.UpdateFromDto(*tt.dto)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.dto.Symbol, s.Symbol)
			}
		})
	}
}

func TestSymbol_ToGetDto(t *testing.T) {
	symbol := &Symbol{
		Symbol:     "BTCUSDT",
		ExchangeID: uuid.New(),
		IsActive:   SymbolActive,
	}
	symbol.ID = uuid.New()

	dto := symbol.ToGetDto()
	assert.Equal(t, symbol.ID.String(), dto.ID)
	assert.Equal(t, symbol.Symbol, dto.Symbol)
	assert.Equal(t, symbol.ExchangeID.String(), dto.ExchangeID)
	assert.Equal(t, symbol.IsActive, dto.IsActive)
}

func TestSymbol_ToDtoUpdate(t *testing.T) {
	symbol := &Symbol{
		Symbol:     "BTCUSDT",
		ExchangeID: uuid.New(),
	}
	symbol.ID = uuid.New()

	dto := symbol.ToDtoUpdate()
	assert.Equal(t, symbol.ID.String(), dto.ID)
	assert.Equal(t, symbol.Symbol, dto.Symbol)
	assert.Equal(t, symbol.ExchangeID.String(), dto.ExchangeID)
}
