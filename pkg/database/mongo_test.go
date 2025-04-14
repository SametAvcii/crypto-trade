package database

import (
	"context"
	"testing"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestBuildMongoURI(t *testing.T) {
	cfg := config.Mongo{
		User:     "testuser",
		Pass:     "testpass",
		Host:     "localhost",
		Port:     "27017",
		Database: "testdb",
	}

	expected := "mongodb://testuser:testpass@localhost:27017/testdb?authSource=admin"
	actual := buildMongoURI(cfg)

	assert.Equal(t, expected, actual)
}

func TestMongoClient(t *testing.T) {
	// Test when client is nil
	client := MongoClient()
	assert.Nil(t, client)

	// Test after initialization
	cfg := config.Mongo{
		User:     "testuser",
		Pass:     "testpass",
		Host:     "localhost",
		Port:     "27017",
		Database: "testdb",
	}
	InitMongo(cfg)

	client = MongoClient()
	assert.NotNil(t, client)
}

func TestConnectToMongo(t *testing.T) {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://invalid:27017")

	err := connectToMongo(ctx, clientOptions)
	assert.Error(t, err)
}

func TestCheckMongoAlive(t *testing.T) {
	cfg := config.Mongo{
		User:     "testuser",
		Pass:     "testpass",
		Host:     "localhost",
		Port:     "27017",
		Database: "testdb",
	}

	done := make(chan bool)
	go func() {
		time.Sleep(20 * time.Second)
		done <- true
	}()

	go CheckMongoAlive(cfg)

	<-done
	assert.False(t, MongoAlive)
}
