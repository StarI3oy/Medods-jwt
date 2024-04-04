package db

import (
	"Star13oy/medods/internal/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// ConnectMongo Подключение к Mongo
func ConnectMongo() (*mongo.Client, error) {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(config.MONGODB_HOST).SetTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return client, err
}
