from dataclasses import asdict, dataclass

import jwt

from observer.common.types import Identifier
from observer.schemas.crypto import PrivateKey


@dataclass
class TokenData:
    user_id: Identifier


class JWTHandler:
    def __init__(self, private_key: PrivateKey):
        self.private_key = private_key

    async def encode(self, payload: TokenData) -> str:
        return jwt.encode(
            asdict(payload),
            self.private_key.key,
            algorithm="RS256",
            headers={
                "kid": self.private_key.hash,
            },
        )

    async def decode(self, token: str) -> TokenData:
        decoded = jwt.decode(token, self.private_key.key.public_key(), algorithms=["RS256"])
        return TokenData(**decoded)
