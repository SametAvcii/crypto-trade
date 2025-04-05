package ctlog

import (
	"github.com/SametAvcii/crypto-trade/pkg/database"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
)

func CreateLog(payload *entities.Log) {
	db := database.PgClient()
	db.Model(&entities.Log{}).Create(payload)
}
