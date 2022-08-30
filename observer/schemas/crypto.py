from dataclasses import dataclass
from enum import Enum

from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey


@dataclass
class PrivateKey:
    hash: str
    key: RSAPrivateKey


class KeyLoaderTypes(str, Enum):
    not_set = "not-set"
    fs = "fs"
    s3 = "s3"
