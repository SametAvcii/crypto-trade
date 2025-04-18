package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SametAvcii/crypto-trade/pkg/config"
	"github.com/SametAvcii/crypto-trade/pkg/consts"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	MongoAlive  bool
)

func InitMongo(cfg config.Mongo) {

	const (
		maxRetries    = consts.MaxRetries
		retryInterval = consts.RetryDelay
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := buildMongoURI(cfg)
	clientOptions := options.Client().ApplyURI(uri)

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err = connectToMongo(ctx, clientOptions); err == nil {
			break
		}

		log.Printf("Failed to connect to MongoDB (attempt %d/%d): %v", attempt, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("MongoDB connection failed after %d attempts: %v", maxRetries, err)
	}

	log.Println("MongoDB connection initialized successfully.")

}

func buildMongoURI(cfg config.Mongo) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func connectToMongo(ctx context.Context, clientOptions *options.ClientOptions) error {
	var err error
	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping to ensure connection is valid
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB ping failed: %v", err)
	}

	return nil
}

func MongoClient() *mongo.Client {
	if mongoClient == nil {
		log.Println("MongoDB is not initialized. Call InitMongo first.")
		return nil
	}
	return mongoClient
}

func MongoDatabase(cfg config.Mongo) *mongo.Database {
	return MongoClient().Database(cfg.Database)
}

func CheckMongoAlive(ctx context.Context, cfg config.Mongo) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if mongoClient == nil {
			log.Println("MongoDB client is nil. Trying to establish a connection...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			uri := buildMongoURI(cfg)
			clientOptions := options.Client().ApplyURI(uri)

			err := connectToMongo(ctx, clientOptions)
			if err != nil {
				MongoAlive = false
				log.Println("Failed to establish MongoDB connection:", err)
				continue
			}
		} else {
			MongoAlive = true
		}

		log.Println("MongoDB connection alive:", MongoAlive)

	}
}
