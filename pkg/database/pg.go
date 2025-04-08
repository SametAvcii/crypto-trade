package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	err         error
	client_once sync.Once
	Alive       bool
)

func InitDB(cfg config.Database) {
	client_once.Do(func() {
		const (
			maxRetries    = 5
			retryInterval = 5 * time.Second
		)

		var (
			dsn = buildDSN(cfg)
			err error
		)

		for attempt := 1; attempt <= maxRetries; attempt++ {
			db, err = openDB(dsn)
			if err == nil {
				break
			}

			log.Printf("Failed to connect to the database (attempt %d/%d): %v", attempt, maxRetries, err)
			time.Sleep(retryInterval)
		}

		if err != nil {
			log.Fatalf("Database connection failed after %d attempts: %v", maxRetries, err)
		}

		if err := runMigrations(); err != nil {
			log.Fatalf("Database migration failed: %v", err)
		}

		log.Println("Database initialized and migrations completed successfully.")
	})
}
func buildDSN(cfg config.Database) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name, cfg.SslMode,
	)
}

func openDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
}

func runMigrations() error {
	return db.AutoMigrate(
		&entities.Log{},
		&entities.Symbol{},
		&entities.Exchange{},
		&entities.SymbolPrice{},
	)
}

func PgClient() *gorm.DB {
	if db == nil {
		log.Panic("Postgres is not initialized. Call InitDB first.")
	}
	return db
}

func CheckPgAlive(dbc config.Database) {
	timeTicker := time.NewTicker(15 * time.Second)

	var (
		dsn = buildDSN(dbc)
		err error
	)

	for range timeTicker.C {
		db, err = openDB(dsn)
		if err != nil {
			Alive = false
		}

		sqldb, err := db.DB()
		if err != nil {
			Alive = false
		}

		err = sqldb.Ping()
		if err != nil {
			Alive = false
		} else {
			Alive = true
		}

		sqldb.Close()

		if !Alive {
			log.Println("Database Alive: ", Alive)
		}
	}
}
