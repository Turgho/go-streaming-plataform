package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDatabase struct {
	Client *mongo.Client
}

func NewMongoDatabase(ctx context.Context, databaseURI, databaseName string) (*MongoDatabase, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(databaseURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo database: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo database: %w", err)
	}

	if err := createMongoIndexes(ctx, client, databaseName); err != nil {
		return nil, fmt.Errorf("failed to create mongo indexes: %w", err)
	}

	return &MongoDatabase{
		Client: client,
	}, nil
}

func createMongoIndexes(ctx context.Context, client *mongo.Client, databaseName string) error {
	collection := client.Database(databaseName).Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	// pega apenas o erro, não pega os indexes criados
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	return nil
}
