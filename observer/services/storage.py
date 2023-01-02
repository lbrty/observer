import os
from datetime import datetime
from pathlib import Path
from typing import IO, Any, List, Protocol

import aiofiles as af
from aiobotocore.session import AioSession, ClientCreatorContext
from aiofiles.os import stat

from observer.common.types import FileInfo, StorageKind
from observer.settings import Settings


class IStorage(Protocol):
    async def ls(self, path: str | Path) -> List[FileInfo]:
        raise NotImplementedError

    async def open(self, path: str | Path) -> IO[Any]:
        raise NotImplementedError

    async def delete(self, path: str):
        raise NotImplementedError


class FSStorage(IStorage):
    async def ls(self, path: str | Path) -> List[FileInfo]:
        items = []
        for root, _, files in os.walk(path):
            for filename in files:
                full_path = os.path.join(root, filename)
                stats = await stat(full_path)
                items.append((datetime.fromtimestamp(stats.st_ctime), full_path))

        return items

    async def open(self, path: str | Path) -> IO[Any]:
        if Path(path).exists():
            fp = await af.open(path, "rb")
            return fp

        raise FileNotFoundError

    async def delete(self, path: str):
        pth = Path(path)
        if pth.is_file():
            pth.unlink(missing_ok=True)


class S3Storage(IStorage):
    def __init__(self, bucket: str, region: str, endpoint_url: str):
        self.bucket = bucket
        self.region = region
        self.endpoint_url = endpoint_url
        self.session = AioSession()

    async def ls(self, path: str | Path) -> List[FileInfo]:
        async with self.s3_client as client:
            result = await client.list_objects_v2(Bucket=self.bucket, Prefix=path)
            return [
                (
                    item["LastModified"],
                    item["Key"],
                )
                for item in result["Contents"]
            ]

    async def open(self, path: str | Path) -> IO[Any]:
        async with self.s3_client as client:
            result = await client.get_object(Bucket=self.bucket, Key=path)
            return result["Body"]

    async def delete(self, path: str):
        async with self.s3_client as client:
            await client.delete_object(Bucket=self.bucket, Key=path)

    @property
    def s3_client(self) -> ClientCreatorContext:
        return self.session.create_client(
            "s3",
            region_name=self.region,
            endpoint_url=self.endpoint_url,
        )


def init_storage(kind: StorageKind, settings: Settings) -> IStorage:
    match kind:
        case StorageKind.fs:
            return FSStorage()
        case StorageKind.s3:
            return S3Storage(settings.s3_bucket, settings.s3_region, settings.s3_endpoint)
        case _ as v:
            raise ValueError(f"Unknown storage type: {v}")
