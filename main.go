package main

import (
	"context"
	"log"
	"os"

	"vault-service/internal/infrastructure"
	"vault-service/internal/service"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Configuration (normally from env vars)
	apiBaseURL := getEnv("API_BASE_URL", "http://localhost:8080")
	s3Bucket := getEnv("S3_BUCKET", "audit-logs")
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")

	// Initialize AWS S3 Client
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)
	storage := infrastructure.NewS3Storage(s3Client, s3Bucket)

	// Initialize MongoDB Client
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	db := infrastructure.NewMongoDB(mongoClient, "audit_vault", "metadata")

	// Initialize API Client
	apiClient := infrastructure.NewAPIClient(apiBaseURL)

	// Initialize VaultService Orchestrator
	vaultService := service.NewVaultService(apiClient, storage, db)

	// For demonstration, we'll run the process once
	log.Println("Starting log processing...")
	err = vaultService.ProcessLogs(ctx)
	if err != nil {
		log.Printf("Error processing logs: %v", err)
	} else {
		log.Println("Successfully processed logs.")
	}

	// Keep alive or exit
	log.Println("Service exiting.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
