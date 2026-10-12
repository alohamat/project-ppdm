package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func ConnectDB(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("falha ao configurar cliente: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("falha no ping com mongodb: %w", err)
	}

	return &Database{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

func (d *Database) Collection(name string) *mongo.Collection {
	return d.DB.Collection(name)
}

func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}