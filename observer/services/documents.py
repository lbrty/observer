from io import BytesIO
from typing import Protocol

from aioboto3 import Session
from aiofile import async_open


class DocumentStoreInterface(Protocol):
    async def load(self, path: str) -> BytesIO:
        raise NotImplementedError


class FS(DocumentStoreInterface):
    async def load(self, path: str) -> BytesIO:
        async with async_open(path, "rb") as fp:
            contents = BytesIO()
            contents.write(await fp.read())
            contents.seek(0)
            return contents


class S3(DocumentStoreInterface):
    def __init__(self, bucket: str):
        self.bucket = bucket

    async def load(self, path: str) -> BytesIO:
        # TODO: implement streaming?
        session = Session()
        async with session.client("s3") as s3:
            obj = await s3.get_object(Bucket=self.bucket, Key=path)
            return obj["Body"]
