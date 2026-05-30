package repository

import (
	"context"
	"fmt"
	"user-service/internal/domain/entities"
	domain "user-service/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName string) domain.UserRepository {
	return &userRepository{
		collection: client.Database(dbName).Collection("users"),
	}
}

// Create implements [repository.UserRepository].
func (u *userRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID implements [repository.UserRepository].
func (u *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User

	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found '%s': %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// GetByEmail implements [repository.UserRepository].
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found '%s': %w", email, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}
