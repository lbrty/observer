import asyncio
import os
import shutil
from datetime import datetime
from pathlib import Path
from typing import Optional, Protocol, Sequence, Union

import aiofiles as af
from aiobotocore.session import AioSession, ClientCreatorContext
from aiofiles.os import stat
from aiofiles.threadpool.binary import AsyncBufferedReader
from botocore.exceptions import ClientError
from starlette import status
from structlog import get_logger

from observer.api.exceptions import InternalError, NotFoundError
from observer.common.types import FileInfo, StorageKind
from observer.settings import Settings

logger = get_logger(name="storage")


class IStorage(Protocol):
    storage_root: Union[str, Path]

    async def ls(self, path: Union[str, Path]) -> Sequence[FileInfo]:
        """List files under the given path

        NOTE: `path` will be prefixed by `storage_root`
        """
        raise NotImplementedError

    async def save(self, path: Union[str, Path], contents: bytes):
        """Save file
        NOTE: `path` will be prefixed by `root`
        """
        raise NotImplementedError

    async def open(self, path: Union[str, Path]) -> AsyncBufferedReader:
        """Open file at `path`
        NOTE: `path` must be absolute path
        """
        raise NotImplementedError

    async def delete(self, path: Union[str, Path]):
        """Delete file
        NOTE: `path` will be prefixed by `root`
        """
        raise NotImplementedError

    async def delete_path(self, path: Union[str, Path]):
        """Delete folder
        NOTE: `path` will be prefixed by `root`
        """
        raise NotImplementedError

    @property
    def root(self) -> str:
        raise NotImplementedError


class FSStorage(IStorage):
    def __init__(self, storage_root: str):
        self.storage_root = storage_root

    async def ls(self, path: Union[str, Path]) -> Sequence[FileInfo]:
        items = []
        for root, _, files in os.walk(path):
            for filename in files:
                full_path = os.path.join(root, filename)
                stats = await stat(full_path)
                items.append((datetime.fromtimestamp(stats.st_ctime), full_path))

        return items

    async def save(self, path: Union[str, Path], contents: bytes):
        pth = Path(self.root) / path
        if not pth.parent.exists():
            pth.parent.mkdir(parents=True, exist_ok=True)

        async with af.open(pth, "wb") as fp:
            await fp.write(contents)

    async def open(self, path: Union[str, Path]) -> AsyncBufferedReader:
        if Path(path).exists():
            fp = await af.open(path, "rb")
            return fp

        raise NotFoundError(message="Document not found")

    async def delete(self, path: Union[str, Path]):
        pth = Path(self.root) / path
        if pth.is_file():
            if pth.exists():
                await af.os.unlink(pth, missing_ok=True)  # type: ignore

    async def delete_path(self, path: Union[str, Path]):
        pth = Path(self.root) / path
        if pth.is_dir():
            shutil.rmtree(pth)

    @property
    def root(self) -> str:
        return self.storage_root


class S3Storage(IStorage):
    def __init__(self, storage_root: str, bucket: str, region: str, endpoint_url: str):
        self.bucket = bucket
        self.region = region
        self.endpoint_url = endpoint_url
        self.storage_root = storage_root
        self.session = AioSession()

    async def ls(self, path: Union[str, Path]) -> Sequence[FileInfo]:
        async with self.s3_client as client:
            result = await client.list_objects_v2(Bucket=self.bucket, Prefix=path)
            if result.get("KeyCount", 0) == 0:
                return []

            return [
                (
                    item["LastModified"],
                    item["Key"],
                )
                for item in result["Contents"]
            ]

    async def save(self, path: Union[str, Path], contents: bytes):
        async with self.s3_client as client:
            try:
                full_path = os.path.join(self.root, path)
                result = await client.put_object(Bucket=self.bucket, Key=full_path, Body=contents)
                if result["ResponseMetadata"]["HTTPStatusCode"] != status.HTTP_200_OK:
                    raise InternalError(message="Unable to upload document to remote storage")
            except ClientError as ex:
                logger.error("Unable to upload document", metadata=ex.response["ResponseMetadata"])
                raise InternalError(message="Unable to upload document")

    async def open(self, path: Union[str, Path]) -> AsyncBufferedReader:
        async with self.s3_client as client:
            try:
                result = await client.get_object(Bucket=self.bucket, Key=path)
                return result["Body"]
            except ClientError as ex:
                if ex.response["Error"]["Code"] == "NoSuchKey":
                    logger.error("Document not found", name=path)
                else:
                    logger.error("Document not found", metadata=ex.response["ResponseMetadata"])

                raise NotFoundError(message="Document not found")

    async def delete(self, path: Union[str, Path]):
        async with self.s3_client as client:
            try:
                if str(path).startswith(str(self.root)):
                    full_path = path
                else:
                    full_path = os.path.join(self.root, path)
                await client.delete_object(Bucket=self.bucket, Key=full_path)
            except ClientError as ex:
                logger.error("Unable to delete document", metadata=ex.response["ResponseMetadata"])
                raise InternalError(message=f"Unable to delete document {full_path}")

    async def delete_path(self, path: Union[str, Path]):
        async with self.s3_client as client:
            try:
                tasks = []
                for _, key in await self.ls(path):
                    tasks.append(self.delete(str(key)))
                await asyncio.gather(*tasks)
                await client.delete_object(Bucket=self.bucket, Key=path)
            except ClientError as ex:
                logger.error("Unable to delete document", metadata=ex.response["ResponseMetadata"])
                raise InternalError(message=f"Unable to delete document {path}")

    @property
    def s3_client(self) -> ClientCreatorContext:
        return self.session.create_client(
            "s3",
            region_name=self.region,
            endpoint_url=self.endpoint_url,
        )

    @property
    def root(self) -> str:
        return self.storage_root


def init_storage(kind: StorageKind, settings: Settings, another_root: Optional[str] = None) -> IStorage:
    """Initialize storage

    Args:
        kind(StorageKind): storage type
        settings(Settings): settings instance
        another_root(Optional[str]): Different storage root atm is used by keychain to load keys.

    Returns: IStorage
    """

    def get_storage(root: str):
        match kind:
            case StorageKind.fs:
                return FSStorage(root)
            case StorageKind.s3:
                if settings.s3_bucket and settings.s3_region and settings.s3_endpoint:
                    return S3Storage(
                        root,
                        settings.s3_bucket,
                        settings.s3_region,
                        settings.s3_endpoint,
                    )

                raise RuntimeError("Please configure S3 storage properly")
            case _ as v:
                raise ValueError(f"Unknown storage type: {v}")

    return get_storage(another_root or settings.storage_root)
