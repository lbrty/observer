import base64
from datetime import datetime, timedelta, timezone
from typing import Protocol, Tuple

from jarowinkler import jarowinkler_similarity
from jwt.exceptions import DecodeError, InvalidAlgorithmError, InvalidSignatureError

from observer.api.exceptions import (
    ForbiddenError,
    InvalidPasswordError,
    NotFoundError,
    PasswordResetCodeExpiredError,
    RegistrationError,
    SimilarPasswordsError,
    TOTPError,
    TOTPRequiredError,
    UnauthorizedError,
    WeakPasswordError,
)
from observer.common import bcrypt
from observer.common.auth import AccessTokenExpirationDelta, RefreshTokenExpirationDelta
from observer.common.bcrypt import is_strong_password
from observer.common.types import Identifier, Role, SomeStr
from observer.entities.users import PasswordReset, User
from observer.schemas.auth import (
    ChangePasswordRequest,
    LoginPayload,
    RegistrationPayload,
    TokenResponse,
)
from observer.schemas.users import NewUserRequest
from observer.services.audit_logs import IAuditService
from observer.services.crypto import ICryptoService
from observer.services.jwt import JWTService, TokenData
from observer.services.mfa import IMFAService
from observer.services.users import IUsersService
from observer.settings import settings


class IAuthService(Protocol):
    crypto_service: ICryptoService
    audits: IAuditService
    mfa_service: IMFAService
    jwt_service: JWTService
    users_service: IUsersService

    async def token_login(self, login_payload: LoginPayload) -> Tuple[User, TokenResponse]:
        raise NotImplementedError

    async def register(self, registration_payload: RegistrationPayload) -> Tuple[User, TokenResponse]:
        raise NotImplementedError

    async def refresh_token(self, refresh_token: str) -> Tuple[TokenData, TokenResponse]:
        raise NotImplementedError

    async def change_password(self, user_id: Identifier, payload: ChangePasswordRequest) -> User:
        raise NotImplementedError

    async def reset_password_request(self, email: str) -> Tuple[User, PasswordReset]:
        raise NotImplementedError

    async def reset_password_with_code(self, code: str, new_password: str) -> User:
        raise NotImplementedError

    async def gen_password(self) -> str:
        raise NotImplementedError

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        raise NotImplementedError

    @property
    def access_token_expiration(self) -> datetime:
        raise NotImplementedError

    @property
    def refresh_token_expiration(self) -> datetime:
        raise NotImplementedError


class AuthService(IAuthService):
    def __init__(
        self,
        crypto_service: ICryptoService,
        mfa_service: IMFAService,
        jwt_service: JWTService,
        users_service: IUsersService,
    ):
        self.crypto_service = crypto_service
        self.jwt_service = jwt_service
        self.mfa_service = mfa_service
        self.users_service = users_service

    async def token_login(self, login_payload: LoginPayload) -> Tuple[User, TokenResponse]:
        if user := await self.users_service.get_by_email(login_payload.email):
            if not user.is_active:
                raise ForbiddenError

            await self.check_totp(user, login_payload.totp_code)
            if bcrypt.check_password(login_payload.password.get_secret_value(), user.password_hash):
                token_response = await self.create_token(user.ref_id)
                return user, token_response

        raise UnauthorizedError(message="Wrong email or password")

    async def register(self, registration_payload: RegistrationPayload) -> Tuple[User, TokenResponse]:
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

        token_response = await self.create_token(user.ref_id)
        return user, token_response

    async def refresh_token(self, refresh_token: str) -> Tuple[TokenData, TokenResponse]:
        try:
            token_data, _ = await self.jwt_service.decode(refresh_token)
            token_response = await self.create_token(token_data.ref_id)
            return token_data, token_response
        except (DecodeError, InvalidAlgorithmError, InvalidSignatureError):
            raise ForbiddenError(message="Invalid refresh token")

    async def change_password(self, user_id: Identifier, payload: ChangePasswordRequest) -> User:
        similarity = jarowinkler_similarity(
            payload.old_password.get_secret_value(),
            payload.new_password.get_secret_value(),
        )
        if similarity >= 0.86:
            raise SimilarPasswordsError(message="Passwords can not be similar")

        if not is_strong_password(payload.new_password.get_secret_value(), settings.password_policy):
            raise WeakPasswordError(message="Given password is weak")

        if user := await self.users_service.get_by_id(user_id):
            await self.check_totp(user, payload.totp_code)
            if bcrypt.check_password(payload.old_password.get_secret_value(), user.password_hash):
                password_hash = bcrypt.hash_password(payload.new_password.get_secret_value())
                return await self.users_service.update_password(user_id, password_hash)
            else:
                raise InvalidPasswordError(message="Invalid password")

        raise UnauthorizedError(message="Unknown user")

    async def reset_password_request(self, email: str) -> Tuple[User, PasswordReset]:
        if user := await self.users_service.get_by_email(email):
            password_reset = await self.users_service.reset_password(user.id)
            return user, password_reset
        else:
            raise UnauthorizedError(message="Unknown user")

    async def reset_password_with_code(self, code: str, new_password: str) -> User:
        if not is_strong_password(new_password, settings.password_policy):
            raise WeakPasswordError(message="Given password is weak")

        if password_reset := await self.users_service.get_password_reset(code):
            now = datetime.now(tz=timezone.utc)
            if now < (password_reset.created_at + timedelta(minutes=settings.password_reset_expiration_minutes)):
                password_hash = bcrypt.hash_password(new_password)
                return await self.users_service.update_password(password_reset.user_id, password_hash)
            else:
                raise PasswordResetCodeExpiredError

        raise NotFoundError

    async def check_totp(self, user: User, totp_code: SomeStr):
        # If MFA is enabled and no TOTP code provided
        # then we need to return HTTP 417 so clients
        # resend auth credentials and TOTP code.
        if user.mfa_enabled and not totp_code:
            raise TOTPRequiredError

        # If MFA is enabled and TOTP given
        # Then we verify it
        if user.mfa_enabled and user.mfa_encrypted_secret:
            # Now we need to decrypt `totp_secret` and verify given `totp_code`
            # if invalid we return `TOTPError`
            keys_hash, encrypted_secret = user.mfa_encrypted_secret.split(":", maxsplit=1)
            decrypted_secret = await self.crypto_service.decrypt(
                keys_hash,
                encrypted_secret.encode(),
            )

            if not totp_code:
                raise TOTPRequiredError

            if not await self.mfa_service.valid(totp_code, decrypted_secret.decode()):
                raise TOTPError(message="Invalid totp code")

    async def gen_password(self) -> str:
        random_bytes = self.crypto_service.gen_key(settings.aes_key_bits)
        return base64.b64encode(random_bytes).decode()

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        payload = TokenData(ref_id=ref_id)
        access_token = await self.jwt_service.encode(payload, self.access_token_expiration)
        refresh_token = await self.jwt_service.encode(payload, self.refresh_token_expiration)
        token = TokenResponse(
            access_token=access_token,
            refresh_token=refresh_token,
        )

        return token

    @property
    def access_token_expiration(self) -> datetime:
        now = datetime.now(tz=timezone.utc)
        return now + AccessTokenExpirationDelta

    @property
    def refresh_token_expiration(self) -> datetime:
        now = datetime.now(tz=timezone.utc)
        return now + RefreshTokenExpirationDelta
