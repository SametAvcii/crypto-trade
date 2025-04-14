package events

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
)

type PgAggTradeHandler struct{}

func (d *PgAggTradeHandler) HandleMessage(msg *sarama.ConsumerMessage) {

	log.Println("Starting Kafka consumer...", msg.Topic)

	var (
		payload     dtos.MongoData
		symbolPrice entities.SymbolPrice
		aggTrade    dtos.AggTrade
	)

	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		log.Printf("Error unmarshalling message for Postgres: %v", err)

	}

	db := database.PgClient()
	json.Unmarshal([]byte(payload.Value), &aggTrade)
	symbolPrice.FromDto(&aggTrade)
	symbolPrice.MongoID = payload.MongoID

	if err := db.Create(&symbolPrice).Error; err != nil {
		log.Printf("Error inserting message into Postgres: %v", err)

	} else {
		//log.Printf("Postgres updated for id %d with value %s", symbolPrice.ID, payload.Value)
	}

	return
}
