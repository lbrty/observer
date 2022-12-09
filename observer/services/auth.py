import base64
from datetime import datetime, timedelta, timezone
from typing import Protocol

from jwt.exceptions import DecodeError, InvalidAlgorithmError, InvalidSignatureError

from observer.api.exceptions import (
    ForbiddenError,
    RegistrationError,
    TOTPError,
    TOTPRequiredError,
    UnauthorizedError,
    WeakPasswordError,
)
from observer.common import bcrypt
from observer.common.bcrypt import is_strong_password
from observer.common.types import Identifier, Role
from observer.schemas.auth import LoginPayload, RegistrationPayload, TokenResponse
from observer.schemas.users import NewUserRequest
from observer.services.crypto import CryptoServiceInterface
from observer.services.jwt import JWTService, TokenData
from observer.services.mfa import MFAServiceInterface
from observer.services.users import UsersServiceInterface
from observer.settings import settings

AccessTokenExpirationMinutes = 15
RefreshTokenExpirationMinutes = 10 * 60 * 24
AccessTokenExpirationDelta = timedelta(minutes=AccessTokenExpirationMinutes)
RefreshTokenExpirationDelta = timedelta(minutes=RefreshTokenExpirationMinutes)


class AuthServiceInterface(Protocol):
    jwt_service: JWTService
    users_service: UsersServiceInterface

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        raise NotImplementedError

    async def register(self, registration_payload: RegistrationPayload) -> TokenResponse:
        raise NotImplementedError

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        raise NotImplementedError

    async def reset_password(self, email: str):
        raise NotImplementedError

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        raise NotImplementedError


class AuthService(AuthServiceInterface):
    def __init__(
        self,
        crypto_service: CryptoServiceInterface,
        mfa_service: MFAServiceInterface,
        jwt_service: JWTService,
        users_service: UsersServiceInterface,
    ):
        self.crypto_service = crypto_service
        self.jwt_service = jwt_service
        self.mfa_service = mfa_service
        self.users_service = users_service

    async def token_login(self, login_payload: LoginPayload) -> TokenResponse:
        if user := await self.users_service.get_by_email(login_payload.email):
            # If MFA is enabled and no TOTP code provided
            # then we need to return HTTP 417 so clients
            # resend auth credentials and TOTP code.
            if user.mfa_enabled and not login_payload.totp_code:
                raise TOTPRequiredError

            # If MFA is enabled and TOTP given
            # Then we verify it
            if user.mfa_enabled:
                # Now we need to decrypt `totp_secret` and verify given `totp_code`
                # if invalid we return `TOTPError`
                keys_hash, encrypted_secret = user.mfa_encrypted_secret.split(":", maxsplit=1)
                decrypted_secret = await self.crypto_service.decrypt(
                    keys_hash, base64.b64decode(encrypted_secret.encode())
                )
                if not await self.mfa_service.valid(login_payload.totp_code, decrypted_secret.decode()):
                    raise TOTPError(message="invalid totp code")

            if bcrypt.check_password(login_payload.password.get_secret_value(), user.password_hash):
                return await self.create_token(user.ref_id)

        raise UnauthorizedError(message="Wrong email or password")

    async def register(self, registration_payload: RegistrationPayload) -> TokenResponse:
        if _ := await self.users_service.get_by_email(registration_payload.email):
            raise RegistrationError(message="User with same e-mail already exists")

        if not is_strong_password(registration_payload.password.get_secret_value(), settings.password_policy):
            raise WeakPasswordError(message="Given password is weak")

        user = await self.users_service.create_user(
            NewUserRequest(
                email=registration_payload.email,
                full_name=None,
                role=Role.guest,
                password=registration_payload.password,
            )
        )

        return await self.create_token(user.ref_id)

    async def refresh_token(self, refresh_token: str) -> TokenResponse:
        try:
            token_data, _ = await self.jwt_service.decode(refresh_token)
            return await self.create_token(token_data.ref_id)
        except (DecodeError, InvalidAlgorithmError, InvalidSignatureError):
            raise ForbiddenError(message="Invalid refresh token")

    async def reset_password(self, email: str):
        ...

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        payload = TokenData(ref_id=ref_id)
        now = datetime.now(tz=timezone.utc)
        access_token = await self.jwt_service.encode(payload, now + AccessTokenExpirationDelta)
        refresh_token = await self.jwt_service.encode(payload, now + RefreshTokenExpirationDelta)
        token = TokenResponse(
            access_token=access_token,
            refresh_token=refresh_token,
        )

        return token
