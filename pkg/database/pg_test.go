package database

import (
	"testing"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildDSN(t *testing.T) {
	cfg := config.Database{
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		Pass:    "password",
		Name:    "testdb",
		SslMode: "disable",
	}

	expected := "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable"
	result := buildDSN(cfg)

	assert.Equal(t, expected, result)
}

func TestInitDB(t *testing.T) {
	cfg := config.Database{
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		Pass:    "password",
		Name:    "testdb",
		SslMode: "disable",
	}

	err := InitDB(cfg)
	assert.NoError(t, err)

	// Test invalid config
	invalidCfg := config.Database{
		Host: "invalid-host",
	}
	err = InitDB(invalidCfg)
	assert.Error(t, err)
}

func TestPgClient(t *testing.T) {
	// Test when db is nil
	client := PgClient()
	assert.Nil(t, client)

	// Initialize DB first
	cfg := config.Database{
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		Pass:    "password",
		Name:    "testdb",
		SslMode: "disable",
	}

	err := InitDB(cfg)
	assert.NoError(t, err)

	// Test after initialization
	client = PgClient()
	assert.NotNil(t, client)
}

func TestCheckPgAlive(t *testing.T) {
	cfg := config.Database{
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		Pass:    "password",
		Name:    "testdb",
		SslMode: "disable",
	}

	done := make(chan bool)
	go func() {
		time.Sleep(20 * time.Second)
		done <- true
	}()

	go CheckPgAlive(cfg)

	<-done
	assert.True(t, Alive)
}
