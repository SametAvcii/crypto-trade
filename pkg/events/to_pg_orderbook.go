package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/cache"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"github.com/SametAvcii/crypto-trade/pkg/ctlog"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/gorm"
)

type PgOrderBookHandler struct{}

func (d *PgOrderBookHandler) HandleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Received message from topic %s", msg.Topic)
	var payload dtos.OrderBook
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		log.Printf("Error unmarshalling message for Postgres: %v", err)
	}

	if err := UpdateOrderBookData(payload.Symbol, payload.Bids, payload.Asks); err != nil {
		log.Printf("Error updating order book data: %v", err)
		ctlog.CreateLog(&entities.Log{
			Title:   "Error updating order book data",
			Message: "Error updating order book data: " + err.Error(),
			Type:    "error",
			Entity:  "order-book",
			Data:    string(msg.Value),
		})
	}
	log.Printf("Order book data updated for symbol %s", payload.Symbol)
	return

}

func UpdateOrderBookData(symbol string, bids, asks [][]string) error {

	rdb := cache.RedisClient()
	bidKey := fmt.Sprintf("order-book-depth:%s:bids", symbol)
	askKey := fmt.Sprintf("order-book-depth:%s:asks", symbol)

	ctx := context.Background()

	oldBids, _ := rdb.HGetAll(ctx, bidKey).Result()
	oldAsks, _ := rdb.HGetAll(ctx, askKey).Result()

	compareAndUpdate("bid", oldBids, bids, symbol)
	compareAndUpdate("ask", oldAsks, asks, symbol)

	newBidMap := make(map[string]string)
	for _, bid := range bids {
		newBidMap[bid[0]] = bid[1]
	}
	rdb.HMSet(ctx, bidKey, newBidMap)

	newAskMap := make(map[string]string)
	for _, ask := range asks {
		newAskMap[ask[0]] = ask[1]
	}
	rdb.HMSet(ctx, askKey, newAskMap)

	return nil
}

func compareAndUpdate(side string, oldData map[string]string, newData [][]string, symbol string) {
	newMap := make(map[string]string)
	for _, entry := range newData {
		newMap[entry[0]] = entry[1]
	}

	for price := range oldData {
		newAmount, exists := newMap[price]
		if !exists || newAmount == "0.00000000" {
			UpdateStatusInDB(symbol, price, side, consts.ClosedOrder)
		}
	}

	for _, entry := range newData {
		price, amount := entry[0], entry[1]
		if amount != "0.00000000" {
			UpsertToDB(symbol, price, amount, side, consts.ActiveOrder)
		}
	}
}

func UpdateStatusInDB(symbol, price, side, status string) {
	db := database.PgClient()
	db.Model(&entities.OrderBook{}).
		Where("symbol = ? AND price = ? AND side = ?", symbol, price, side).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

}

func UpsertToDB(symbol, price, amount, side, status string) {

	db := database.PgClient()
	var existing entities.OrderBook
	err := db.Where("symbol = ? AND price = ? AND side = ?", symbol, price, side).First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Kayıt bulunamadı, yeni kayıt oluştur
			newOrder := entities.OrderBook{
				Symbol: symbol,
				Price:  price,
				Amount: amount,
				Side:   side,
				Status: status,
			}
			db.Create(&newOrder)
		} else {
			// Başka bir hata varsa logla veya yönet
			log.Printf("Error querying order: %v", err)
		}
	} else {
		// Kayıt bulundu, güncelle
		existing.Amount = amount
		existing.Status = status
		db.Save(&existing)
	}
}
