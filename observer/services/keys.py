import hashlib
import os
from glob import glob
from typing import List, Protocol

from cryptography.hazmat.primitives.serialization import load_pem_private_key
from structlog import get_logger

from observer.schemas.crypto import KeyLoaderTypes, PrivateKey

logger = get_logger(service="keys")


class Keychain(Protocol):
    keys: List[PrivateKey] = []
    name: KeyLoaderTypes = KeyLoaderTypes.not_set

    async def load(self, path: str):
        raise NotImplementedError

    async def find(self, key_hash: str) -> PrivateKey | None:
        for key in self.keys:
            if key.hash == key_hash:
                return key
        return None


class FS(Keychain):
    name = KeyLoaderTypes.fs

    async def load(self, path: str):
        keys = []
        if key_list := glob(f"{path}/*.pem"):
            for key in key_list:
                fp = os.open(key, os.O_RDONLY)
                stats = os.fstat(fp)
                file_bytes = os.read(fp, stats.st_size)
                h = hashlib.new("sha256", file_bytes)
                pem_private_key = load_pem_private_key(file_bytes, password=None)
                private_key = PrivateKey(
                    hash=str(h.hexdigest())[:16].upper(),
                    private_key=pem_private_key,  # type:ignore
                )
                keys.append((stats.st_ctime, private_key))

                logger.info("Loaded RSAPrivateKey", SHA256=private_key.hash)
                os.close(fp)

        sorted(keys, key=lambda item: item[0])
        self.keys = [item[1] for item in keys]


class S3(Keychain):
    name = KeyLoaderTypes.s3

    async def load(self, path: str):
        pass


class UnknownKeyLoaderError(ValueError):
    pass


def get_key_loader(loader_type: KeyLoaderTypes) -> Keychain:
    match loader_type:
        case KeyLoaderTypes.fs:
            return FS()
        case KeyLoaderTypes.s3:
            return S3()
        case _:
            raise UnknownKeyLoaderError
