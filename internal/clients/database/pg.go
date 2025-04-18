package database

import (
	"context"
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
	dsn         string
)

func InitDB(cfg config.Database) error {
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
		return fmt.Errorf("Database connection failed after %d attempts: %w", maxRetries, err)
	}

	sqldb, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed to get DB from GORM: %w", err)
	}

	sqldb.SetMaxIdleConns(3)
	sqldb.SetMaxOpenConns(90)
	sqldb.SetConnMaxLifetime(time.Hour)
	sqldb.SetConnMaxIdleTime(2 * time.Second)

	if err := runMigrations(); err != nil {
		return fmt.Errorf("Database migration failed: %w", err)
	}

	Seed()

	log.Println("Database initialized and migrations completed successfully.")
	return nil
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
		&entities.SignalInterval{},
		&entities.Signal{},
		&entities.Candlestick{},
		&entities.OrderBook{},
	)
}

func PgClient() *gorm.DB {
	if db == nil {
		log.Println("Postgres is not initialized. Call InitDB first.")
		return nil
	}

	// Bağlantıyı test et
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get underlying DB:", err)
		err := InitDB(config.ReadValue().Database)
		if err != nil {
			log.Println("Failed to reinitialize database:", err)
			return nil
		}
		return db
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("PostgreSQL connection lost. Reconnecting...")
		err := InitDB(config.ReadValue().Database)
		if err != nil {
			log.Println("Failed to reinitialize database:", err)
			return nil
		}
		log.Println("PostgreSQL reconnected successfully.")

	}

	return db
}
func CheckPgAlive(ctx context.Context, dbc config.Database) {
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
		log.Println("Postgres connection alive: ", Alive)
	}
}
