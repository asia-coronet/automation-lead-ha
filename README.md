# Audit Vault Service - System Under Test (SUT)

This project implements the "Audit Vault Service", a data pipeline service designed to be tested for integrity and resilience.

## Architecture

The service follows a clean architecture with a clear separation of concerns:

- **Orchestrator (`vault_service/orchestrator.py`)**: Manages the high-level data flow (Extract -> Load -> Audit).
- **Interfaces (`vault_service/interfaces.py`)**: Defines abstract contracts for external dependencies using `abc.ABC`.
- **Infrastructure (`vault_service/infrastructure.py`)**: Concrete implementations for API fetching, S3 storage, and MongoDB persistence.

### Dependency Injection Decisions

To ensure testability and extensibility:
- **Constructor Injection**: All external clients (API, Storage, Database) are injected via the `VaultService` constructor.
- **Abstract Base Classes**: The service depends on interfaces, not concrete implementations, allowing for easy mocking or swapping of infrastructure.
- **Fixture Support**: Concrete implementations are instantiated at the entry point (`main.py`) or within test fixtures, keeping the core logic isolated.

## Execution

### Run with Docker

```bash
docker build -t vault-automation-test .
docker run --rm vault-automation-test
```

### Data Flow Logic

1. **Extract**: Fetches log data from an external REST API (with retry logic).
2. **Load (S3)**: Stores raw JSON payload in AWS S3.
3. **Audit (DB)**: Persists metadata (S3 key, timestamp, status) in MongoDB.
   - If S3 upload fails, the MongoDB audit record is **not** written, ensuring pipeline integrity.
