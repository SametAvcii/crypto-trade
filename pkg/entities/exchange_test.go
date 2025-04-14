package entities

import (
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/stretchr/testify/assert"
)

func TestExchange_FromDto(t *testing.T) {
	exchange := &Exchange{}
	req := dtos.AddExchangeReq{
		Name:  "Binance",
		WsUrl: "wss://ws-api.binance.com",
	}

	exchange.FromDto(req)

	assert.Equal(t, req.Name, exchange.Name)
	assert.Equal(t, req.WsUrl, exchange.WsUrl)
	assert.Equal(t, ExchangeActive, exchange.IsActive)
}

func TestExchange_ToDto(t *testing.T) {
	exchange := &Exchange{
		Name:  "Binance",
		WsUrl: "wss://ws-api.binance.com",
	}

	res := exchange.ToDto()

	assert.Equal(t, exchange.ID.String(), res.ID)
	assert.Equal(t, exchange.Name, res.Name)
	assert.Equal(t, exchange.WsUrl, res.WsUrl)
}

func TestExchange_FromDtoUpdate(t *testing.T) {
	exchange := &Exchange{}
	req := dtos.UpdateExchangeReq{
		Name:  "Kucoin",
		WsUrl: "wss://ws-api.kucoin.com",
	}

	exchange.FromDtoUpdate(req)

	assert.Equal(t, req.Name, exchange.Name)
	assert.Equal(t, req.WsUrl, exchange.WsUrl)
}

func TestExchange_ToDtoUpdate(t *testing.T) {
	exchange := &Exchange{
		Name:  "Kucoin",
		WsUrl: "wss://ws-api.kucoin.com",
	}

	res := exchange.ToDtoUpdate()

	assert.Equal(t, exchange.ID.String(), res.ID)
	assert.Equal(t, exchange.Name, res.Name)
	assert.Equal(t, exchange.WsUrl, res.WsUrl)
}

func TestExchange_ToDtoGet(t *testing.T) {
	exchange := &Exchange{
		Name:  "Binance",
		WsUrl: "wss://ws-api.binance.com",
	}

	res := exchange.ToDtoGet()

	assert.Equal(t, exchange.ID.String(), res.ID)
	assert.Equal(t, exchange.Name, res.Name)
	assert.Equal(t, exchange.WsUrl, res.WsUrl)
}