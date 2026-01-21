package service

import (
	"context"
	"fmt"
	"time"

	"vault-service/internal/domain"
)

// VaultService is the orchestrator that manages the data pipeline flow.
type VaultService struct {
	apiClient     domain.APIClient
	objectStorage domain.ObjectStorage
	database      domain.Database
}

// NewVaultService creates a new instance of VaultService with the provided dependencies.
func NewVaultService(api domain.APIClient, storage domain.ObjectStorage, db domain.Database) *VaultService {
	return &VaultService{
		apiClient:     api,
		objectStorage: storage,
		database:      db,
	}
}

// ProcessLogs executes the end-to-end data pipeline: Fetch -> Upload -> Audit.
func (s *VaultService) ProcessLogs(ctx context.Context) error {
	// 1. Extract - Fetch log data from REST API
	payload, err := s.apiClient.FetchLogEntry(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch log entry: %w", err)
	}

	// Generate S3 key (e.g., using timestamp)
	timestamp := time.Now()
	s3Key := fmt.Sprintf("logs/%d.json", timestamp.UnixNano())

	// 2. Load (S3) - Store raw JSON payload in S3
	err = s.objectStorage.Upload(ctx, s3Key, payload)
	if err != nil {
		return fmt.Errorf("failed to upload payload to S3: %w", err)
	}

	// 3. Audit (DB) - Persist metadata record in MongoDB
	metadata := domain.Metadata{
		S3Key:     s3Key,
		Timestamp: timestamp,
		Status:    "archived",
	}

	err = s.database.SaveMetadata(ctx, metadata)
	if err != nil {
		return fmt.Errorf("failed to save metadata to database: %w", err)
	}

	return nil
}
