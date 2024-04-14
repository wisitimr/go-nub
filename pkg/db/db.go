package db

import (
	"context"
	"nub/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client instance
func Connect(ctx context.Context, config config.Database) (*mongo.Database, error) {
	option := options.Client().ApplyURI(config.URI).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	option.SetMinPoolSize(10)
	option.SetMaxPoolSize(100)
	option.SetMaxConnIdleTime(2 * time.Second)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return nil, err
	}
	return client.Database(config.Name), err
}
