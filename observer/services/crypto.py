import hashlib
from dataclasses import dataclass
from enum import Enum
from glob import glob
from typing import List, Protocol

from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey
from cryptography.hazmat.primitives.serialization import load_pem_private_key

from observer.settings import settings


@dataclass
class PrivateKey:
    hash: str
    key: RSAPrivateKey


class KeyLoaderTypes(str, Enum):
    not_set = "not-set"
    fs = "fs"
    s3 = "s3"


class KeychainLoader(Protocol):
    keys: List[PrivateKey] = []
    name: KeyLoaderTypes = KeyLoaderTypes.not_set

    async def load(self, path: str):
        raise NotImplementedError


class FSLoader(KeychainLoader):
    name = KeyLoaderTypes.fs

    async def load(self, path: str):
        if key_list := glob(f"{str(settings.key_store_path)}/**/*.pem"):
            for key in key_list:
                with open(key, "rb") as fp:
                    file_bytes = fp.read()
                    h = hashlib.new("sha256", file_bytes)
                    self.keys.append(
                        PrivateKey(
                            hash=str(h.hexdigest())[:16].upper(),
                            key=load_pem_private_key(
                                file_bytes,
                                password=None,
                            ),
                        )
                    )


class S3Loader(KeychainLoader):
    name = KeyLoaderTypes.s3

    async def load(self, path: str):
        pass
