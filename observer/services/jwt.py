from dataclasses import asdict, dataclass
from datetime import datetime, timezone

import jwt

from observer.common.types import Identifier
from observer.schemas.crypto import PrivateKey


@dataclass
class TokenData:
    ref_id: Identifier


class JWTService:
    def __init__(self, private_key: PrivateKey):
        self.private_key = private_key

    async def encode(self, payload: TokenData, expiration: datetime) -> str:
        """JWT encode `TokenData` with given expiration datetime.
        Here we use RS256 and private key to generate JWT token
        and with this we also include the following token header info

        1. kid,
        2. iat,
        3. exp

        Args:
            payload(TokenData): JWT payload
            expiration(datetime): expiration datetime

        Returns:
            str: encoded token
        """
        return jwt.encode(
            asdict(payload),
            self.private_key.private_key,
            algorithm="RS256",
            headers={
                "kid": self.private_key.hash,
                "iat": datetime.now(tz=timezone.utc).timestamp(),
                "exp": expiration.timestamp(),
            },
        )

    async def decode(self, token: str) -> tuple[TokenData, dict]:
        """Decode JWT token
        We decode and verify token at the same time so errors listed
        below have to be handled by callers.

        1. `DecodeError`,
        2. `InvalidAlgorithmError`,
        3. `InvalidSignatureError`.

        NOTE: we don't give leeway to any token

        Args:
            token(str): encoded JWT token

        Returns:
            tuple[TokenData, dict]: a pair of `TokenData` and jwt token headers dictionary
        """
        decoded = jwt.decode(token, self.private_key.private_key.public_key(), algorithms=["RS256"])
        return TokenData(**decoded), jwt.get_unverified_header(token)
