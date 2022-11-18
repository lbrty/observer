from datetime import datetime, timedelta, timezone
from typing import Protocol

from observer.api.exceptions import UnauthorizedError
from observer.common import bcrypt
from observer.schemas.auth import LoginPayload, TokenResponse
from observer.services.jwt import JWTHandler, TokenData
from observer.services.users import UsersServiceInterface

AccessTokenExpirationMinutes = 15
RefreshTokenExpirationMinutes = 10 * 60 * 24
AccessTokenExpirationDelta = timedelta(minutes=AccessTokenExpirationMinutes)
RefreshTokenExpirationDelta = timedelta(minutes=RefreshTokenExpirationMinutes)


class AuthServiceInterface(Protocol):
    jwt_handler: JWTHandler
    users_service: UsersServiceInterface

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        raise NotImplementedError()

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        raise NotImplementedError()

    async def reset_password(self, email: str):
        raise NotImplementedError()


class AuthService(AuthServiceInterface):
    def __init__(self, jwt_handler: JWTHandler, users_service: UsersServiceInterface):
        self.jwt_handler = jwt_handler
        self.users_service = users_service

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        if user := await self.users_service.get_by_email(login_payload.email):
            if bcrypt.check_password(login_payload.password.get_secret_value(), user.password_hash):
                payload = TokenData(ref_id=user.ref_id)
                now = datetime.now(tz=timezone.utc)
                return TokenResponse(
                    access_token=self.jwt_handler.encode(payload, now + AccessTokenExpirationDelta),
                    refresh_token=self.jwt_handler.encode(payload, now + RefreshTokenExpirationDelta),
                )

        raise UnauthorizedError(message="Wrong email or password")

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        ...

    async def reset_password(self, email: str):
        ...
