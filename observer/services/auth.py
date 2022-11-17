from typing import Protocol

from observer.schemas.auth import LoginPayload, TokenResponse
from observer.services.jwt import JWTHandler
from observer.services.users import UsersServiceInterface


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
        ...

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        ...

    async def reset_password(self, email: str):
        ...
