from io import BytesIO
from typing import Protocol

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
    async def load(self, path: str) -> BytesIO:
        raise NotImplementedError
