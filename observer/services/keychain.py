import fnmatch
import hashlib
from typing import List, Protocol

from cryptography.hazmat.primitives.serialization import load_pem_private_key
from structlog import get_logger

from observer.api.exceptions import InternalError
from observer.schemas.crypto import PrivateKey
from observer.services.storage import IStorage

logger = get_logger(service="keys")


class IKeychain(Protocol):
    keys: List[PrivateKey] = []

    async def load(self, path: str, storage: IStorage):
        raise NotImplementedError

    async def find(self, key_hash: str) -> PrivateKey:
        raise NotImplementedError


class Keychain(IKeychain):
    async def load(self, path: str, storage: IStorage):
        keys = []
        files = await storage.ls(path)
        for creation_time, filename in files:
            if not fnmatch.fnmatch(str(filename), "*.pem"):
                continue

            fp = await storage.open(filename)
            file_bytes = await fp.read()
            private_key = await self.parse_key(file_bytes, str(filename))
            keys.append((creation_time, private_key))

            logger.info("Loaded RSAPrivateKey", SHA256=private_key.hash)

        sorted(keys, key=lambda item: item[0])
        self.keys = [item[1] for item in keys]

    async def find(self, key_hash: str) -> PrivateKey:
        for key in self.keys:
            if key.hash == key_hash:
                return key

        raise InternalError(message=f"Private key with hash={key_hash} not found")

    async def parse_key(self, key_bytes: bytes, filename: str) -> PrivateKey:
        h = hashlib.new("sha256", key_bytes)
        pem_private_key = load_pem_private_key(key_bytes, password=None)
        private_key = PrivateKey(
            filename=filename,
            hash=str(h.hexdigest())[:16].upper(),
            private_key=pem_private_key,  # type:ignore
        )
        return private_key
