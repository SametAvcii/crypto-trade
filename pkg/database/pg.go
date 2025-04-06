package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	err         error
	client_once sync.Once
)

func InitDB(dbc config.Database) {
	client_once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Istanbul", dbc.Host, dbc.Port, dbc.User, dbc.Pass, dbc.Name)
		db, err = gorm.Open(
			postgres.New(
				postgres.Config{
					DSN:                  dsn,
					PreferSimpleProtocol: true, // daha az kaynak
				},
			),
		)
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(
			&entities.Log{},
			&entities.Symbol{},
			&entities.Exchange{},
			&entities.SymbolPrice{},
		)

	})
}

func PgClient() *gorm.DB {
	if db == nil {
		log.Panic("Postgres is not initialized. Call InitDB first.")
	}
	return db
}
