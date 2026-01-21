import time
import logging
from datetime import datetime
from typing import Optional
from requests.exceptions import RequestException
from .interfaces import APIClient, ObjectStorage, Database

logger = logging.getLogger(__name__)

class VaultService:
    def __init__(self, api_client: APIClient, storage: ObjectStorage, database: Database):
        self.api_client = api_client
        self.storage = storage
        self.database = database

    def process_logs(self, max_retries: int = 3) -> None:
        """
        Extract -> Load (S3) -> Audit (DB)
        """
        payload = self._fetch_with_retry(max_retries)
        if not payload:
            logger.error("Failed to fetch payload from API after retries.")
            return

        s3_key = f"logs/{int(time.time() * 1000)}.json"

        try:
            # Load to S3
            self.storage.upload_payload(s3_key, payload)
        except Exception as e:
            logger.error(f"Failed to upload to S3: {e}")
            # Scenario 2: No MongoDB audit record is written on storage failure
            return

        # Audit to DB
        metadata = {
            "s3_object_key": s3_key,
            "timestamp": datetime.utcnow().isoformat(),
            "status": "archived"
        }
        self.database.write_metadata(metadata)
        logger.info("Successfully processed and audited log entry.")

    def _fetch_with_retry(self, max_retries: int) -> Optional[dict]:
        for attempt in range(max_retries):
            try:
                return self.api_client.fetch_log_entry()
            except RequestException as e:
                logger.warning(f"API attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(0.1)  # Brief sleep before retry
                continue
        return None
