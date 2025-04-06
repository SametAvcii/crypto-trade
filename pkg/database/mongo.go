package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoOnce   sync.Once
	mongoErr    error
)

func InitMongo(cfg config.Mongo) {
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)

		clientOptions := options.Client().ApplyURI(uri)

		mongoClient, mongoErr = mongo.Connect(ctx, clientOptions)
		if mongoErr != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", mongoErr)
		}

		// Ping to ensure connection is valid
		if err := mongoClient.Ping(ctx, nil); err != nil {
			log.Fatalf("MongoDB ping failed: %v", err)
		}

		log.Println("MongoDB connection initialized successfully.")
	})
}

func MongoClient() *mongo.Client {
	if mongoClient == nil {
		log.Panic("MongoDB is not initialized. Call InitMongo first.")
	}
	return mongoClient
}

func MongoDatabase(cfg config.Mongo) *mongo.Database {
	return MongoClient().Database(cfg.Database)
}
