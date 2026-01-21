import json
import requests
import boto3
from pymongo import MongoClient
from typing import Dict, Any
from .interfaces import APIClient, ObjectStorage, Database

class RESTAPIClient(APIClient):
    def __init__(self, base_url: str):
        self.base_url = base_url

    def fetch_log_entry(self) -> Dict[str, Any]:
        response = requests.get(f"{self.base_url}/logs")
        response.raise_for_status()
        return response.json()

class S3Storage(ObjectStorage):
    def __init__(self, bucket_name: str):
        self.bucket_name = bucket_name
        self.s3 = boto3.client("s3")

    def upload_payload(self, key: str, payload: Dict[str, Any]) -> None:
        self.s3.put_object(
            Bucket=self.bucket_name,
            Key=key,
            Body=json.dumps(payload)
        )

class MongoDB(Database):
    def __init__(self, connection_string: str, db_name: str, collection_name: str):
        self.client = MongoClient(connection_string)
        self.db = self.client[db_name]
        self.collection = self.db[collection_name]

    def write_metadata(self, metadata: Dict[str, Any]) -> None:
        self.collection.insert_one(metadata)
