import os
import logging
from vault_service.infrastructure import RESTAPIClient, S3Storage, MongoDB
from vault_service.orchestrator import VaultService

logging.basicConfig(level=logging.INFO)

def main():
    # Configuration from environment variables
    api_url = os.getenv("API_BASE_URL", "http://localhost:8080")
    s3_bucket = os.getenv("S3_BUCKET", "audit-logs")
    mongo_uri = os.getenv("MONGO_URI", "mongodb://localhost:27017")

    # Dependency Injection
    api_client = RESTAPIClient(base_url=api_url)
    storage = S3Storage(bucket_name=s3_bucket)
    database = MongoDB(
        connection_string=mongo_uri,
        db_name="audit_vault",
        collection_name="metadata"
    )

    service = VaultService(api_client, storage, database)

    print("Starting VaultService Log Processing...")
    service.process_logs()
    print("Log processing cycle completed.")

if __name__ == "__main__":
    main()
