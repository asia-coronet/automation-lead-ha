package domain

import (
	"context"
	"time"
)

// Metadata represents the audit record stored in the database.
type Metadata struct {
	S3Key     string    `json:"s3_key"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// APIClient defines the contract for fetching log entries from an external API.
type APIClient interface {
	FetchLogEntry(ctx context.Context) ([]byte, error)
}

// ObjectStorage defines the contract for storing raw log payloads.
type ObjectStorage interface {
	Upload(ctx context.Context, key string, data []byte) error
}

// Database defines the contract for persisting audit metadata.
type Database interface {
	SaveMetadata(ctx context.Context, metadata Metadata) error
}
