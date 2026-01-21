from abc import ABC, abstractmethod
from typing import Dict, Any

class APIClient(ABC):
    @abstractmethod
    def fetch_log_entry(self) -> Dict[str, Any]:
        """Fetch log entry from external REST API."""
        pass

class ObjectStorage(ABC):
    @abstractmethod
    def upload_payload(self, key: str, payload: Dict[str, Any]) -> None:
        """Upload raw JSON payload to S3."""
        pass

class Database(ABC):
    @abstractmethod
    def write_metadata(self, metadata: Dict[str, Any]) -> None:
        """Write audit metadata into MongoDB."""
        pass
