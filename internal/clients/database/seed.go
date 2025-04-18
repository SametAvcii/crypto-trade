package database

import (
	"log"

	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
)

func Seed() {
	log.Println("Seeding database...")
	db := PgClient()
	exchange := &entities.Exchange{
		Name:     consts.Binance,
		WsUrl:    "wss://stream.binance.com:443/ws",
		RestUrl:  "https://api.binance.com/api/v3",
		IsActive: consts.Active,
	}

	err := db.Where("name = ?", exchange.Name).First(exchange).Error
	if err != nil {
		err = db.Create(exchange).Error
		if err != nil {
			errLog := &entities.Log{
				Title:   "Error creating exchange from seed",
				Message: "Error creating exchange from seed: " + err.Error(),
				Entity:  "exchange",
				Type:    "error",
			}
			db.Model(&entities.Log{}).Create(errLog)

			log.Println("Error creating exchange:", err)
			return
		}
	}

	var symbol = &entities.Symbol{
		Symbol:     "btcusdt",
		ExchangeID: exchange.ID,
		IsActive:   consts.Active,
	}
	err = db.Where("symbol = ? and exchange_id = ?", symbol.Symbol, exchange.ID).First(symbol).Error
	if err != nil {
		err = db.Create(symbol).Error
		if err != nil {
			errLog := &entities.Log{
				Title:   "Error creating symbol from seed",
				Message: "Error creating symbol from seed: " + err.Error(),
				Entity:  "symbol",
				Type:    "error",
			}
			db.Model(&entities.Log{}).Create(errLog)

			log.Println("Error creating symbol:", err)
			return
		}
	}

	var signalInterval = &entities.SignalInterval{
		Symbol:     symbol.Symbol,
		Interval:   "1m",
		ExchangeID: exchange.ID,
		IsActive:   consts.Active,
	}
	err = db.Where("symbol = ? and interval = ? and exchange_id = ?", signalInterval.Symbol, signalInterval.Interval, exchange.ID).First(signalInterval).Error
	if err != nil {
		err = db.Create(signalInterval).Error
		if err != nil {
			errLog := &entities.Log{
				Title:   "Error creating signal interval from seed",
				Message: "Error creating signal interval from seed: " + err.Error(),
				Entity:  "signal_interval",
				Type:    "error",
			}

			db.Model(&entities.Log{}).Create(errLog)
			log.Println("Error creating signal interval:", err)
			return
		}
	}

	infoLog := &entities.Log{
		Title:   "Dummy data created",
		Message: "Dummy data created successfully",
		Entity:  "seed",
		Type:    "info",
	}
	db.Model(&entities.Log{}).Create(infoLog)

	log.Println("Dummy data created successfully")

}
