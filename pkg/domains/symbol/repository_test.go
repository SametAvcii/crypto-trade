package symbol_test

import (
	"context"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/domains/symbol"
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

	err = db.AutoMigrate(&entities.Symbol{})
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

func TestAddSymbol(t *testing.T) {
	db := setupTestDB(t)
	repo := symbol.NewRepo(db)

	req := dtos.AddSymbolReq{
		Symbol: "BTCUSDT",
	}

	res, err := repo.AddSymbol(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Symbol, res.Symbol)
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := symbol.NewRepo(db)

	addRes, err := repo.AddSymbol(context.Background(), dtos.AddSymbolReq{
		Symbol: "BTCUSDT",
	})
	assert.NoError(t, err)

	res, err := repo.GetByID(context.Background(), addRes.ID)
	assert.NoError(t, err)
	assert.Equal(t, addRes.Symbol, res.Symbol)
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := symbol.NewRepo(db)

	_, err := repo.AddSymbol(context.Background(), dtos.AddSymbolReq{Symbol: "BTCUSDT"})
	assert.NoError(t, err)

	_, err = repo.AddSymbol(context.Background(), dtos.AddSymbolReq{Symbol: "ETHUSDT"})
	assert.NoError(t, err)

	res, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := symbol.NewRepo(db)

	addRes, err := repo.AddSymbol(context.Background(), dtos.AddSymbolReq{Symbol: "BTCUSDT"})
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), addRes.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), addRes.ID)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := symbol.NewRepo(db)

	addRes, err := repo.AddSymbol(context.Background(), dtos.AddSymbolReq{Symbol: "BTCUSDT"})
	assert.NoError(t, err)

	updateReq := dtos.UpdateSymbolReq{
		ID:     addRes.ID,
		Symbol: "ETHUSDT",
	}

	res, err := repo.Update(context.Background(), updateReq)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Symbol, res.Symbol)
}
