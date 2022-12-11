import base64
from datetime import datetime, timedelta, timezone
from typing import Protocol, Tuple

from jwt.exceptions import DecodeError, InvalidAlgorithmError, InvalidSignatureError
from starlette.background import BackgroundTasks

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
from observer.entities.users import User
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.auth import LoginPayload, RegistrationPayload, TokenResponse
from observer.schemas.users import NewUserRequest
from observer.services.audit_logs import AuditServiceInterface
from observer.services.crypto import CryptoServiceInterface
from observer.services.jwt import JWTService, TokenData
from observer.services.mailer import EmailMessage, MailerInterface
from observer.services.mfa import MFAServiceInterface
from observer.services.users import UsersServiceInterface
from observer.settings import settings

AccessTokenExpirationMinutes = 15
RefreshTokenExpirationMinutes = 10 * 60 * 24
AccessTokenExpirationDelta = timedelta(minutes=AccessTokenExpirationMinutes)
RefreshTokenExpirationDelta = timedelta(minutes=RefreshTokenExpirationMinutes)


class AuthServiceInterface(Protocol):
    tag: str
    jwt_service: JWTService
    users_service: UsersServiceInterface

    async def token_login(self, login_payload: LoginPayload) -> Tuple[User, TokenResponse]:
        raise NotImplementedError

    async def register(self, registration_payload: RegistrationPayload) -> TokenResponse:
        raise NotImplementedError

    async def refresh_token(self, refresh_token: str) -> Tuple[TokenData, TokenResponse]:
        raise NotImplementedError

    async def reset_password(self, email: str, metadata: dict = None):
        raise NotImplementedError

    async def create_token(self, ref_id: Identifier) -> TokenResponse:
        raise NotImplementedError

    async def create_log(self, ref: str, expires_in: timedelta, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError


class AuthService(AuthServiceInterface):
    tag: str = "origin=auth,source=service:auth"

    def __init__(
        self,
        crypto_service: CryptoServiceInterface,
        audits: AuditServiceInterface,
        mailer: MailerInterface,
        mfa_service: MFAServiceInterface,
        jwt_service: JWTService,
        users_service: UsersServiceInterface,
    ):
        self.crypto_service = crypto_service
        self.audits = audits
        self.mailer = mailer
        self.jwt_service = jwt_service
        self.mfa_service = mfa_service
        self.users_service = users_service
        self.tasks = BackgroundTasks()

    async def token_login(self, login_payload: LoginPayload) -> Tuple[User, TokenResponse]:
        if user := await self.users_service.get_by_email(login_payload.email):
            await self.check_totp(user, login_payload.totp_code)
            if bcrypt.check_password(login_payload.password.get_secret_value(), user.password_hash):
                token_response = await self.create_token(user.ref_id)
                return user, token_response

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

        token_response = await self.create_token(user.ref_id)
        self.tasks.add_task(
            self.audits.add_event,
            NewAuditLog(
                ref=f"{self.tag},action=token:register",
                data=dict(ref_id=user.ref_id, role=user.role.value),
                expires_at=None,
            ),
        )
        return token_response

    async def refresh_token(self, refresh_token: str) -> Tuple[TokenData, TokenResponse]:
        try:
            token_data, _ = await self.jwt_service.decode(refresh_token)
            token_response = await self.create_token(token_data.ref_id)
            return token_data, token_response
        except (DecodeError, InvalidAlgorithmError, InvalidSignatureError):
            raise ForbiddenError(message="Invalid refresh token")

    async def reset_password(self, email: str, metadata: dict = None):
        now = datetime.now(tz=timezone.utc)
        data = dict(email=email)
        if metadata:
            data = {
                **data,
                **metadata,
            }

        if user := await self.users_service.get_by_email(email):
            password_reset = await self.users_service.reset_password(user.id)
            reset_link = f"{settings.app_domain}/{settings.password_reset_url.format(code=password_reset.code)}"
            self.tasks.add_task(
                self.mailer.send,
                EmailMessage(
                    to_email=user.email,
                    from_email=settings.from_email,
                    subject=settings.mfa_reset_subject,
                    body=f"To reset you password please use the following link {reset_link}.",
                ),
            )
            self.tasks.add_task(
                self.audits.add_event,
                NewAuditLog(
                    ref=f"{self.tag},action=reset:password",
                    data=data,
                    expires_at=now + timedelta(days=settings.auth_audit_event_lifetime_days),
                ),
            )

    async def check_totp(self, user: User, totp_code: str | None):
        # If MFA is enabled and no TOTP code provided
        # then we need to return HTTP 417 so clients
        # resend auth credentials and TOTP code.
        if user.mfa_enabled and not totp_code:
            raise TOTPRequiredError

        # If MFA is enabled and TOTP given
        # Then we verify it
        if user.mfa_enabled:
            # Now we need to decrypt `totp_secret` and verify given `totp_code`
            # if invalid we return `TOTPError`
            keys_hash, encrypted_secret = user.mfa_encrypted_secret.split(":", maxsplit=1)
            decrypted_secret = await self.crypto_service.decrypt(
                keys_hash,
                base64.b64decode(encrypted_secret.encode()),
            )
            if not await self.mfa_service.valid(totp_code, decrypted_secret.decode()):
                raise TOTPError(message="Invalid totp code")

    async def create_log(self, ref: str, expires_in: timedelta, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=now + expires_in,
        )

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
