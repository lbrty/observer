from dataclasses import dataclass

from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey


@dataclass
class PrivateKey:
    filename: str
    hash: str
    private_key: RSAPrivateKey
