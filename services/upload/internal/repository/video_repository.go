package repository

import (
	"context"
	"fmt"
	"time"
	"upload-service/internal/domain/entities"
	domain "upload-service/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type videoRepository struct {
	collection *mongo.Collection
}

func NewVideoRepository(client *mongo.Client, dbName string) domain.VideoRepository {
	return &videoRepository{
		collection: client.Database(dbName).Collection("videos"),
	}
}

// Create implements [repositories.VideoRepository].
func (v *videoRepository) Create(ctx context.Context, video *entities.Video) error {
	_, err := v.collection.InsertOne(ctx, video)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}
	return nil
}

// GetByID implements [repositories.VideoRepository].
func (v *videoRepository) GetByID(ctx context.Context, id string) (*entities.Video, error) {
	var video entities.Video

	err := v.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&video)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("video not found '%s': %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find video: %w", err)
	}
	return &video, nil
}

// ListByUserID implements [repositories.VideoRepository].
func (v *videoRepository) ListByUserID(ctx context.Context, userID string) ([]*entities.Video, error) {
	cursor, err := v.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to find videos: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []*entities.Video
	if err := cursor.All(ctx, &videos); err != nil {
		return nil, fmt.Errorf("failed to get videos: %w", err)
	}

	return videos, nil
}

// UpdateStatus implements [repositories.VideoRepository].
func (v *videoRepository) UpdateStatus(ctx context.Context, id string, status entities.VideoStatus) error {
	_, err := v.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update video status: %w", err)
	}
	return nil
}

// UpdateMetadata implements [repositories.VideoRepository].
func (v *videoRepository) UpdateMetadata(ctx context.Context, id string, duration float64, size int64) error {
	_, err := v.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"duration":   duration,
			"size":       size,
			"updated_at": time.Now().UTC(),
		}},
	)
	if err != nil {
		return fmt.Errorf("failed to update video metadata: %w", err)
	}
	return nil
}
