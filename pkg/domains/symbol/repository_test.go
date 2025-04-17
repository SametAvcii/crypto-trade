package symbol_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SametAvcii/crypto-trade/pkg/domains/symbol"
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
		&entities.Symbol{},
	)

	return gormDB, mock
}

func TestAddSymbolRepo(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	req := dtos.AddSymbolReq{
		Symbol:     "btcusdt",
		ExchangeID: "540e8400-e29b-41d4-a716-446655440000",
	}

	// Mock the DB interaction for adding a new symbol
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "symbols"`).
		WithArgs(
			sqlmock.AnyArg(), // ID (uuid)
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			"btcusdt",
			"540e8400-e29b-41d4-a716-446655440000",
			1,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := repo.AddSymbol(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "btcusdt", res.Symbol)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	rows := sqlmock.NewRows([]string{"id", "symbol", "created_at", "updated_at"}).
		AddRow("550e8400-e29b-41d4-a716-446655440000", "BTCUSDT", time.Now(), time.Now())

	query := regexp.QuoteMeta(`SELECT * FROM "symbols" WHERE id = $1`)
	mock.ExpectQuery(query).WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).WillReturnRows(rows)

	res, err := repo.GetByID(context.Background(), "550e8400-e29b-41d4-a716-446655440000")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", res.Symbol)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	rows := sqlmock.NewRows([]string{"id", "symbol", "created_at", "updated_at", "exchange_id", "status"}).
		AddRow("550e8400-e29b-41d4-a716-446655440000", "BTCUSDT", time.Now(), time.Now(), "540e8400-e29b-41d4-a716-446655440000", 1).
		AddRow("650e8400-e29b-41d4-a716-446655440000", "ETHUSDT", time.Now(), time.Now(), "540e8400-e29b-41d4-a716-446655440000", 1)

	mock.ExpectQuery(`SELECT \* FROM "symbols"`).
		WillReturnRows(rows)

	res, err := repo.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	query := regexp.QuoteMeta(`UPDATE "symbols" SET "deleted_at"=$1 WHERE id = $2 AND "symbols"."deleted_at" IS NULL`)

	mock.ExpectBegin()
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg(), "550e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), "550e8400-e29b-41d4-a716-446655440000")

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	rows := sqlmock.NewRows([]string{"id", "symbol", "created_at", "updated_at"}).
		AddRow("550e8400-e29b-41d4-a716-446655440000", "BTCUSDT", "2021-01-01", "2021-01-01")

	mock.ExpectQuery(`SELECT \* FROM "symbols" WHERE id = \$1`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").
		WillReturnRows(rows)

	req := dtos.UpdateSymbolReq{
		ID:     "550e8400-e29b-41d4-a716-446655440000",
		Symbol: "ETHUSDT",
	}
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "symbols" SET`).
		WithArgs("ethusdt", sqlmock.AnyArg(), sqlmock.AnyArg(), "550e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := repo.Update(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, req.Symbol, res.Symbol)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateSymbol(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := symbol.NewRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "symbols" WHERE id = $1 AND "symbols"."deleted_at" IS NULL ORDER BY "symbols"."id" LIMIT $2`)).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "symbol", "exchange_id"}).
			AddRow("550e8400-e29b-41d4-a716-446655440000", "ethusdt", "650e8400-e29b-41d4-a716-446655440000"))

	// Mock the UPDATE query
	req := dtos.UpdateSymbolReq{
		ID:         "550e8400-e29b-41d4-a716-446655440000",
		Symbol:     "ethusdt",
		ExchangeID: "650e8400-e29b-41d4-a716-446655440000",
	}

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE "symbols" SET`).
		WithArgs(sqlmock.AnyArg(), "ethusdt", req.ExchangeID, req.ID, req.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	res, err := repo.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "ethusdt", res.Symbol)

	assert.NoError(t, mock.ExpectationsWereMet())
}
