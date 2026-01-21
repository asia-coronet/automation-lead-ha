package infrastructure

import (
	"context"

	"vault-service/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB(client *mongo.Client, dbName, collectionName string) *MongoDB {
	return &MongoDB{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (m *MongoDB) SaveMetadata(ctx context.Context, metadata domain.Metadata) error {
	_, err := m.collection.InsertOne(ctx, metadata)
	return err
}
