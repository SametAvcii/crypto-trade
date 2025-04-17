package signal_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SametAvcii/crypto-trade/pkg/domains/signal"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	gormDB.AutoMigrate(
		&entities.SignalInterval{},
	)

	return gormDB, mock
}

func TestAddSignalIntervals(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := signal.NewRepo(db)

	req := dtos.AddSignalIntervalReq{
		Interval:   "1m",
		Symbol:     "BTCUSDT",
		ExchangeId: "550e8400-e29b-41d4-a716-446655440000",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "signal_intervals"`).
		WithArgs(
			sqlmock.AnyArg(), // ID (uuid)
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			"btcusdt",
			"1m",
			"550e8400-e29b-41d4-a716-446655440000", // ExchangeID (uuid)
			1,                                      // IsActive
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := repo.AddSignalIntervals(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, req.Interval, res.Interval)
	assert.Equal(t, "btcusdt", res.Symbol)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSignalInterval(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := signal.NewRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "signal_intervals" WHERE id = $1 AND "signal_intervals"."deleted_at" IS NULL ORDER BY "signal_intervals"."id" LIMIT $2`)).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "interval", "symbol", "exchange_id"}).
			AddRow("550e8400-e29b-41d4-a716-446655440000", "1m", "ethusdt", "650e8400-e29b-41d4-a716-446655440000"))

	// Mock the UPDATE query
	req := dtos.UpdateSignalIntervalReq{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		Interval:   "5m",
		Symbol:     "ETHUSDT",
		ExchangeId: "650e8400-e29b-41d4-a716-446655440000",
	}

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "signal_intervals" SET`).
		WithArgs(sqlmock.AnyArg(), "ethusdt", "5m", req.ExchangeId, req.ID, req.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	res, err := repo.UpdateSignalInterval(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, req.Interval, res.Interval)
	assert.Equal(t, "ethusdt", res.Symbol)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSignalInterval(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := signal.NewRepo(db)

	mock.ExpectBegin()
	query := regexp.QuoteMeta(`UPDATE "signal_intervals" SET "deleted_at"=$1 WHERE id = $2 AND "signal_intervals"."deleted_at" IS NULL`)
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg(), "750e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteSignalInterval(context.Background(), "750e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetSignalInterval(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := signal.NewRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "interval", "symbol", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"1m",
		"BTCUSDT",
		time.Now(),
		time.Now(),
		nil, // assuming "deleted_at" is nil (not deleted)
	)
	query := regexp.QuoteMeta(`SELECT * FROM "signal_intervals" WHERE id = $1 AND "signal_intervals"."deleted_at" IS NULL ORDER BY "signal_intervals"."id" LIMIT $2`)
	mock.ExpectQuery(query).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
		WillReturnRows(rows)

	res, err := repo.GetSignalInterval(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)
	assert.Equal(t, "1m", res.Interval)
	assert.Equal(t, "BTCUSDT", res.Symbol)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllSignalIntervals_Repo(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := signal.NewRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "symbol", "interval", "exchange_id", "is_active",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"Interval 1",
		"1h",
		"aaaabbbb-cccc-dddd-eeee-ffffffffffff",
		1,
	).AddRow(
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"Interval 2",
		"4h",
		"11112222-3333-4444-5555-666677778888",
		2,
	)

	mock.ExpectQuery(`SELECT \* FROM "signal_intervals"`).
		WillReturnRows(rows)

	got, err := repo.GetAllSignalIntervals(context.Background())
	assert.NoError(t, err)

	want := []dtos.GetSignalIntervalRes{
		{
			ID:         "550e8400-e29b-41d4-a716-446655440000",
			Symbol:     "Interval 1",
			Interval:   "1h",
			ExchangeId: "aaaabbbb-cccc-dddd-eeee-ffffffffffff",
			IsActive:   1,
		},
		{
			ID:         "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			Symbol:     "Interval 2",
			Interval:   "4h",
			ExchangeId: "11112222-3333-4444-5555-666677778888",
			IsActive:   2,
		},
	}
	assert.Equal(t, want, got)

	assert.NoError(t, mock.ExpectationsWereMet())
}
