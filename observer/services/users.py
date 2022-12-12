import base64
from datetime import datetime, timedelta, timezone
from typing import Protocol

import shortuuid

from observer.api.exceptions import TOTPError
from observer.common import bcrypt
from observer.common.types import Identifier
from observer.entities.base import SomeUser
from observer.entities.users import (
    Confirmation,
    NewUser,
    PasswordReset,
    User,
    UserUpdate,
)
from observer.repositories.users import UsersRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.users import (
    NewUserRequest,
    UserMFAUpdateRequest,
    UserResponse,
    UsersResponse,
)
from observer.services.crypto import CryptoServiceInterface


class UsersServiceInterface(Protocol):
    tag: str
    repo: UsersRepositoryInterface

    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> SomeUser:
        raise NotImplementedError

    async def create_user(self, new_user: NewUserRequest) -> User:
        raise NotImplementedError

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        raise NotImplementedError

    async def update_mfa(self, user_id: Identifier, updates: UserMFAUpdateRequest):
        raise NotImplementedError

    async def reset_mfa(self, user_id: Identifier):
        raise NotImplementedError

    async def check_backup_code(self, user_backup_codes: str, given_backup_code: str):
        raise NotImplementedError

    async def create_confirmation(self, user_id: Identifier, code: str) -> Confirmation:
        raise NotImplementedError

    async def get_confirmation(self, code: str) -> Confirmation:
        raise NotImplementedError

    async def reset_password(self, user_id: Identifier) -> PasswordReset:
        raise NotImplementedError

    async def get_password_reset(self, code: str) -> PasswordReset | None:
        raise NotImplementedError

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        raise NotImplementedError


class UsersService(UsersServiceInterface):
    tag: str = "origin=mfa,source=service:user"

    def __init__(self, users_repository: UsersRepositoryInterface, crypto_service: CryptoServiceInterface):
        self.repo = users_repository
        self.crypto_service = crypto_service

    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        return await self.repo.get_by_id(user_id)

    async def get_by_ref_id(self, ref_id: Identifier) -> SomeUser:
        return await self.repo.get_by_ref_id(ref_id)

    async def get_by_email(self, email: str) -> SomeUser:
        return await self.repo.get_by_email(email)

    async def create_user(self, new_user: NewUserRequest) -> User:
        ref_id = shortuuid.uuid(name=new_user.email)
        password_hash = bcrypt.hash_password(new_user.password.get_secret_value())
        user = NewUser(
            ref_id=ref_id,
            email=new_user.email,
            full_name=new_user.full_name,
            password_hash=password_hash,
            role=new_user.role,
            is_active=True,
            is_confirmed=False,
        )
        return await self.repo.create_user(user)

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        return await self.repo.update_password(user_id, new_password_hash)

    async def update_mfa(self, user_id: Identifier, updates: UserMFAUpdateRequest):
        user_update = UserUpdate(
            mfa_enabled=updates.mfa_enabled,
            mfa_encrypted_secret=updates.mfa_encrypted_secret,
            mfa_encrypted_backup_codes=updates.mfa_encrypted_backup_codes,
        )
        await self.repo.update_user(user_id, user_update)

    async def reset_mfa(self, user_id: Identifier):
        await self.repo.reset_mfa(user_id)

    async def check_backup_code(self, user_backup_codes: str, given_backup_code: str):
        keys_hash, encrypted_backup_codes = user_backup_codes.split(":", maxsplit=1)
        decrypted_backup_codes = await self.crypto_service.decrypt(
            keys_hash, base64.b64decode(encrypted_backup_codes.encode())
        )
        if given_backup_code not in decrypted_backup_codes.decode().split(","):
            raise TOTPError(message="invalid backup code")

    async def create_confirmation(self, user_id: Identifier, code: str) -> Confirmation:
        return await self.repo.create_confirmation(user_id, shortuuid.uuid())

    async def get_confirmation(self, code: str) -> Confirmation | None:
        return await self.repo.get_confirmation(code)

    async def reset_password(self, user_id: Identifier) -> PasswordReset:
        return await self.repo.create_password_reset_code(user_id, shortuuid.uuid())

    async def get_password_reset(self, code: str) -> PasswordReset | None:
        return await self.repo.get_password_reset(code)

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if expires_in:
            expires_at = now + expires_in

        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=expires_at,
        )

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        return UserResponse(**user.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        return UsersResponse(
            total=total,
            items=[UserResponse(**user.dict()) for user in user_list],
        )
