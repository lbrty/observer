import hashlib
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


class FS(Keychain):
    name = KeyLoaderTypes.fs

    async def load(self, path: str):
        if key_list := glob(f"{path}/*.pem"):
            for key in key_list:
                with open(key, "rb") as fp:
                    file_bytes = fp.read()
                    h = hashlib.new("sha256", file_bytes)
                    private_key = PrivateKey(
                        hash=str(h.hexdigest())[:16].upper(),
                        private_key=load_pem_private_key(
                            file_bytes,
                            password=None,
                        ),
                    )
                    self.keys.append(private_key)

                    logger.info("Loaded RSAPrivateKey", SHA256=h)


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
