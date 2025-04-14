package exchange_test

import (
	"context"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/domains/exchange"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&entities.Exchange{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Errorf("failed to get database instance: %v", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	return db
}

func TestAddExchange(t *testing.T) {
	db := setupTestDB(t)
	repo := exchange.NewRepo(db)

	req := dtos.AddExchangeReq{
		Name:  "Binance",
		WsUrl: "wss://stream.binance.com:9443/ws",
	}

	res, err := repo.AddExchange(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, res.Name)
}

func TestUpdateExchange(t *testing.T) {
	db := setupTestDB(t)
	repo := exchange.NewRepo(db)

	// Create initial data
	addRes, err := repo.AddExchange(context.Background(), dtos.AddExchangeReq{
		Name:  "Binance",
		WsUrl: "wss://stream.binance.com:9443/ws",
	})
	assert.NoError(t, err)

	req := dtos.UpdateExchangeReq{
		ID:    addRes.ID,
		Name:  "Updated Binance",
		WsUrl: "wss://stream.binance.com:9443/ws",
	}

	res, err := repo.UpdateExchange(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Binance", res.Name)
}

func TestDeleteExchange(t *testing.T) {
	db := setupTestDB(t)
	repo := exchange.NewRepo(db)
	addRes, err := repo.AddExchange(context.Background(), dtos.AddExchangeReq{
		Name:  "ToDelete",
		WsUrl: "wss://stream.binance.com:9443/ws",
	})
	assert.NoError(t, err)

	// Verify exchange exists before deletion
	_, err = repo.GetExchangeById(context.Background(), addRes.ID)
	assert.NoError(t, err)

	err = repo.DeleteExchange(context.Background(), addRes.ID)
	assert.NoError(t, err)

	// Ensure it is gone
	_, err = repo.GetExchangeById(context.Background(), addRes.ID)
	assert.Error(t, err)
}

func TestGetExchangeById(t *testing.T) {
	db := setupTestDB(t)
	repo := exchange.NewRepo(db)
	addRes, err := repo.AddExchange(context.Background(), dtos.AddExchangeReq{
		Name:  "FetchMe",
		WsUrl: "test.com",
	})
	assert.NoError(t, err)

	res, err := repo.GetExchangeById(context.Background(), addRes.ID)
	assert.NoError(t, err)
	assert.Equal(t, addRes.Name, res.Name)
}

func TestGetAllExchanges(t *testing.T) {
	db := setupTestDB(t)
	repo := exchange.NewRepo(db)

	// Add two exchanges
	var err error
	_, err = repo.AddExchange(context.Background(), dtos.AddExchangeReq{Name: "A", WsUrl: "wss://binance.com"})
	assert.NoError(t, err)
	_, err = repo.AddExchange(context.Background(), dtos.AddExchangeReq{Name: "B", WsUrl: "wss://kraken.com"})
	assert.NoError(t, err)

	res, err := repo.GetAllExchanges(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
