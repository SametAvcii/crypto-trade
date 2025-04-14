package signal_test

import (
	"context"
	"testing"

	"github.com/SametAvcii/crypto-trade/pkg/domains/signal"
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

	err = db.AutoMigrate(&entities.SignalInterval{})
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

func TestAddSignalIntervals(t *testing.T) {
	db := setupTestDB(t)
	repo := signal.NewRepo(db)

	req := dtos.AddSignalIntervalReq{
		Interval: "1m",
		Symbol:   "BTCUSDT",
	}

	res, err := repo.AddSignalIntervals(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Interval, res.Interval)
}

func TestUpdateSignalInterval(t *testing.T) {
	db := setupTestDB(t)
	repo := signal.NewRepo(db)

	addRes, err := repo.AddSignalIntervals(context.Background(), dtos.AddSignalIntervalReq{
		Interval: "1m",
		Symbol:   "BTCUSDT",
	})
	assert.NoError(t, err)

	req := dtos.UpdateSignalIntervalReq{
		ID:       addRes.ID,
		Interval: "5m",
		Symbol:   "ETHUSDT",
	}

	res, err := repo.UpdateSignalInterval(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "5m", res.Interval)
}

func TestDeleteSignalInterval(t *testing.T) {
	db := setupTestDB(t)
	repo := signal.NewRepo(db)

	addRes, err := repo.AddSignalIntervals(context.Background(), dtos.AddSignalIntervalReq{
		Interval: "1m",
		Symbol:   "BTCUSDT",
	})
	assert.NoError(t, err)

	err = repo.DeleteSignalInterval(context.Background(), addRes.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetSignalInterval(context.Background(), addRes.ID)
	assert.Error(t, err)
}

func TestGetSignalInterval(t *testing.T) {
	db := setupTestDB(t)
	repo := signal.NewRepo(db)

	addRes, err := repo.AddSignalIntervals(context.Background(), dtos.AddSignalIntervalReq{
		Interval: "1m",
		Symbol:   "BTCUSDT",
	})
	assert.NoError(t, err)

	res, err := repo.GetSignalInterval(context.Background(), addRes.ID)
	assert.NoError(t, err)
	assert.Equal(t, addRes.Interval, res.Interval)
}

func TestGetAllSignalIntervals(t *testing.T) {
	db := setupTestDB(t)
	repo := signal.NewRepo(db)

	_, err := repo.AddSignalIntervals(context.Background(), dtos.AddSignalIntervalReq{
		Interval: "1m",
		Symbol:   "BTCUSDT",
	})
	assert.NoError(t, err)

	_, err = repo.AddSignalIntervals(context.Background(), dtos.AddSignalIntervalReq{
		Interval: "5m",
		Symbol:   "ETHUSDT",
	})
	assert.NoError(t, err)

	res, err := repo.GetAllSignalIntervals(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}
