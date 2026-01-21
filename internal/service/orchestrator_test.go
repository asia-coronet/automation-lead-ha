package service

import (
	"context"
	"errors"
	"testing"
	"vault-service/internal/domain"
)

// MockAPIClient is a mock implementation of APIClient
type MockAPIClient struct {
	FetchFunc func(ctx context.Context) ([]byte, error)
}

func (m *MockAPIClient) FetchLogEntry(ctx context.Context) ([]byte, error) {
	return m.FetchFunc(ctx)
}

// MockObjectStorage is a mock implementation of ObjectStorage
type MockObjectStorage struct {
	UploadFunc func(ctx context.Context, key string, data []byte) error
}

func (m *MockObjectStorage) Upload(ctx context.Context, key string, data []byte) error {
	return m.UploadFunc(ctx, key, data)
}

// MockDatabase is a mock implementation of Database
type MockDatabase struct {
	SaveFunc func(ctx context.Context, metadata domain.Metadata) error
}

func (m *MockDatabase) SaveMetadata(ctx context.Context, metadata domain.Metadata) error {
	return m.SaveFunc(ctx, metadata)
}

func TestVaultService_ProcessLogs_HappyPath(t *testing.T) {
	ctx := context.Background()

	mockAPI := &MockAPIClient{
		FetchFunc: func(ctx context.Context) ([]byte, error) {
			return []byte(`{"log": "data"}`), nil
		},
	}

	mockStorage := &MockObjectStorage{
		UploadFunc: func(ctx context.Context, key string, data []byte) error {
			return nil
		},
	}

	mockDB := &MockDatabase{
		SaveFunc: func(ctx context.Context, metadata domain.Metadata) error {
			if metadata.Status != "archived" {
				t.Errorf("expected status 'archived', got %s", metadata.Status)
			}
			return nil
		},
	}

	svc := NewVaultService(mockAPI, mockStorage, mockDB)
	err := svc.ProcessLogs(ctx)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVaultService_ProcessLogs_StorageFailure(t *testing.T) {
	ctx := context.Background()

	mockAPI := &MockAPIClient{
		FetchFunc: func(ctx context.Context) ([]byte, error) {
			return []byte(`{"log": "data"}`), nil
		},
	}

	mockStorage := &MockObjectStorage{
		UploadFunc: func(ctx context.Context, key string, data []byte) error {
			return errors.New("S3 403 Forbidden")
		},
	}

	mockDB := &MockDatabase{
		SaveFunc: func(ctx context.Context, metadata domain.Metadata) error {
			t.Error("database should not be called on storage failure")
			return nil
		},
	}

	svc := NewVaultService(mockAPI, mockStorage, mockDB)
	err := svc.ProcessLogs(ctx)

	if err == nil {
		t.Fatal("expected error on storage failure, got nil")
	}
}
