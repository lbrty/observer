from datetime import datetime, timedelta, timezone
from typing import Protocol

from jwt.exceptions import DecodeError, InvalidAlgorithmError, InvalidSignatureError

from observer.api.exceptions import ForbiddenError, UnauthorizedError
from observer.common import bcrypt
from observer.common.types import Identifier
from observer.schemas.auth import LoginPayload, TokenResponse
from observer.services.jwt import JWTService, TokenData
from observer.services.users import UsersServiceInterface

AccessTokenExpirationMinutes = 15
RefreshTokenExpirationMinutes = 10 * 60 * 24
AccessTokenExpirationDelta = timedelta(minutes=AccessTokenExpirationMinutes)
RefreshTokenExpirationDelta = timedelta(minutes=RefreshTokenExpirationMinutes)


class AuthServiceInterface(Protocol):
    jwt_handler: JWTService
    users_service: UsersServiceInterface

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        raise NotImplementedError()

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        raise NotImplementedError()

    async def reset_password(self, email: str):
        raise NotImplementedError()

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        raise NotImplementedError()


class AuthService(AuthServiceInterface):
    def __init__(self, jwt_handler: JWTService, users_service: UsersServiceInterface):
        self.jwt_handler = jwt_handler
        self.users_service = users_service

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        if user := await self.users_service.get_by_email(login_payload.email):
            if bcrypt.check_password(login_payload.password.get_secret_value(), user.password_hash):
                return await self.create_token(user.ref_id)

        raise UnauthorizedError(message="Wrong email or password")

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        try:
            token_data, _ = await self.jwt_handler.decode(refresh_token)
            return await self.create_token(token_data.ref_id)
        except (DecodeError, InvalidAlgorithmError, InvalidSignatureError):
            raise ForbiddenError(message="Invalid refresh token")

    async def reset_password(self, email: str):
        ...

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        payload = TokenData(ref_id=ref_id)
        now = datetime.now(tz=timezone.utc)
        return TokenResponse(
            access_token=self.jwt_handler.encode(payload, now + AccessTokenExpirationDelta),
            refresh_token=self.jwt_handler.encode(payload, now + RefreshTokenExpirationDelta),
        )
