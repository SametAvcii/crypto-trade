package exchange_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SametAvcii/crypto-trade/pkg/domains/exchange"
	"github.com/SametAvcii/crypto-trade/pkg/dtos"
	"github.com/SametAvcii/crypto-trade/pkg/entities"
	"github.com/google/uuid"
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
		&entities.Exchange{},
	)

	return gormDB, mock
}

func TestNewRepo(t *testing.T) {
	// Setup mock database
	db, _ := setupMockDB(t)

	repo := exchange.NewRepo(db)

	assert.NotNil(t, repo, "Repository should not be nil")

	_, ok := repo.(exchange.Repository)
	assert.True(t, ok, "Returned value should implement Repository interface")
}

func TestAddExchange(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := exchange.NewRepo(db)

	req := dtos.AddExchangeReq{
		Name:  "Binance",
		WsUrl: "wss://stream.binance.com:9443/ws",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "exchanges"`).
		WithArgs(
			sqlmock.AnyArg(), // ID (uuid)
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			req.Name,
			req.WsUrl,
			"", // RestUrl
			1,  // IsActive
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	res, err := repo.AddExchange(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, res.Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateExchange(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := exchange.NewRepo(db)
	exchangeID := uuid.New().String()
	req := dtos.UpdateExchangeReq{
		ID:    exchangeID,
		Name:  "Updated Exchange",
		WsUrl: "wss://updated.url",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "exchanges" WHERE id = $1 AND "exchanges"."deleted_at" IS NULL ORDER BY "exchanges"."id" LIMIT $2`)).
		WithArgs(exchangeID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "ws_url"}).
			AddRow(exchangeID, "old name", "wss://old.url"))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "exchanges" SET`).
		WithArgs(sqlmock.AnyArg(), req.Name, req.WsUrl, sqlmock.AnyArg(), req.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	res, err := repo.UpdateExchange(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.ID, res.ID)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.WsUrl, res.WsUrl)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExchangeById(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := exchange.NewRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "name", "ws_url",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"FetchMe",
		"test.com")

	query := regexp.QuoteMeta(`SELECT * FROM "exchanges" WHERE id = $1 AND "exchanges"."deleted_at" IS NULL ORDER BY "exchanges"."id" LIMIT $2`)

	mock.ExpectQuery(query).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", 1).
		WillReturnRows(rows)

	res, err := repo.GetExchangeById(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)
	assert.Equal(t, "FetchMe", res.Name)
	assert.Equal(t, "test.com", res.WsUrl)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllExchanges(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := exchange.NewRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "name", "ws_url",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"A",
		"wss://binance.com",
	).AddRow(
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"B",
		"wss://kraken.com",
	)

	mock.ExpectQuery(`SELECT \* FROM "exchanges"`).
		WillReturnRows(rows)

	got, err := repo.GetAllExchanges(context.Background())
	assert.NoError(t, err)

	want := []dtos.GetExchangeRes{
		{
			ID:    "550e8400-e29b-41d4-a716-446655440000",
			Name:  "A",
			WsUrl: "wss://binance.com",
		},
		{
			ID:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			Name:  "B",
			WsUrl: "wss://kraken.com",
		},
	}
	assert.Equal(t, want, got)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteExchange(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := exchange.NewRepo(db)
	mock.ExpectBegin()
	query := regexp.QuoteMeta(`UPDATE "exchanges" SET "deleted_at"=$1 WHERE id = $2 AND "exchanges"."deleted_at" IS NULL`)
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg(), "550e8400-e29b-41d4-a716-446655440000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteExchange(context.Background(), "550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
