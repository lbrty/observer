import os
import shutil
from datetime import datetime
from pathlib import Path
from typing import IO, Any, List, Protocol

import aiofiles as af
from aiobotocore.session import AioSession, ClientCreatorContext
from aiofiles.os import stat
from botocore.exceptions import ClientError
from structlog import get_logger

from observer.api.exceptions import InternalError
from observer.common.types import FileInfo, StorageKind
from observer.settings import Settings

logger = get_logger(name="storage")


class IStorage(Protocol):
    documents_root: str | Path

    async def ls(self, path: str | Path) -> List[FileInfo]:
        """List files under the given path

        NOTE: `path` will be prefixed by `documents_root`
        """
        raise NotImplementedError

    async def save(self, path: str | Path, contents: bytes):
        """Save file
        NOTE: `path` will be prefixed by `documents_root`
        """
        raise NotImplementedError

    async def open(self, path: str | Path) -> IO[Any]:
        """Open file at `path`
        NOTE: `path` must be absolute path
        """
        raise NotImplementedError

    async def delete(self, path: str):
        """Delete file
        NOTE: `path` will be prefixed by `documents_root`
        """
        raise NotImplementedError

    @property
    def root(self) -> str:
        raise NotImplementedError


class FSStorage(IStorage):
    def __init__(self, documents_root: str):
        self.documents_root = documents_root

    async def ls(self, path: str | Path) -> List[FileInfo]:
        items = []
        for root, _, files in os.walk(Path(self.root) / path):
            for filename in files:
                full_path = os.path.join(root, filename)
                stats = await stat(full_path)
                items.append((datetime.fromtimestamp(stats.st_ctime), full_path))

        return items

    async def save(self, path: str | Path, contents: bytes):
        pth = Path(self.root) / path
        if not pth.parent.exists():
            pth.parent.mkdir(parents=True, exist_ok=True)

        async with af.open(pth, "wb") as fp:
            await fp.write(contents)

    async def open(self, path: str | Path) -> IO[Any]:
        if Path(path).exists():
            fp = await af.open(path, "rb")
            return fp

        raise FileNotFoundError

    async def delete(self, path: str):
        pth = Path(self.root) / path
        if pth.is_file():
            pth.unlink(missing_ok=True)
        elif pth.is_dir():
            shutil.rmtree(pth)

    @property
    def root(self) -> str:
        return self.documents_root


class S3Storage(IStorage):
    def __init__(self, documents_root: str, bucket: str, region: str, endpoint_url: str):
        self.bucket = bucket
        self.region = region
        self.endpoint_url = endpoint_url
        self.documents_root = documents_root
        self.session = AioSession()

    async def ls(self, path: str | Path) -> List[FileInfo]:
        async with self.s3_client as client:
            full_path = os.path.join(self.root, path)
            result = await client.list_objects_v2(Bucket=self.bucket, Prefix=full_path)
            return [
                (
                    item["LastModified"],
                    item["Key"],
                )
                for item in result["Contents"]
                if result["KeyCount"] > 0
            ]

    async def save(self, path: str | Path, contents: bytes):
        async with self.s3_client as client:
            try:
                # TODO: check for success
                full_path = os.path.join(self.root, path)
                await client.put_object(Bucket=self.bucket, Key=full_path, Body=contents)
            except ClientError as ex:
                logger.error("Unable to upload document", metadata=ex.response["ResponseMetadata"])
                raise InternalError(message="Unable to upload document")

    async def open(self, path: str | Path) -> IO[Any]:
        async with self.s3_client as client:
            try:
                result = await client.get_object(Bucket=self.bucket, Key=path)
                return result["Body"]
            except ClientError as ex:
                if ex.response["Error"]["Code"] == "NoSuchKey":
                    logger.error("Document not found", name=path)
                else:
                    logger.error("Document not found", metadata=ex.response["ResponseMetadata"])
                    raise InternalError(message="Document not found")

    async def delete(self, path: str):
        async with self.s3_client as client:
            try:
                full_path = os.path.join(self.root, path)
                await client.delete_object(Bucket=self.bucket, Key=full_path)
            except ClientError as ex:
                logger.error("Unable to delete document", metadata=ex.response["ResponseMetadata"])
                raise InternalError(message=f"Unable to delete document {full_path}")

    @property
    def s3_client(self) -> ClientCreatorContext:
        return self.session.create_client(
            "s3",
            region_name=self.region,
            endpoint_url=self.endpoint_url,
        )

    @property
    def root(self) -> str:
        return self.documents_root


def init_storage(kind: StorageKind, settings: Settings) -> IStorage:
    match kind:
        case StorageKind.fs:
            return FSStorage(str(settings.documents_path))
        case StorageKind.s3:
            return S3Storage(
                str(settings.documents_path),
                settings.s3_bucket,
                settings.s3_region,
                settings.s3_endpoint,
            )
        case _ as v:
            raise ValueError(f"Unknown storage type: {v}")
