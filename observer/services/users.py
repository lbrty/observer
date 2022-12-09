import base64
from typing import Protocol

import shortuuid

from observer.api.exceptions import TOTPError
from observer.common import bcrypt
from observer.common.types import Identifier
from observer.entities.base import SomeUser
from observer.entities.users import NewUser, User, UserUpdate
from observer.repositories.users import UsersRepositoryInterface
from observer.schemas.users import (
    NewUserRequest,
    UserMFAUpdateRequest,
    UserResponse,
    UsersResponse,
)
from observer.services.crypto import CryptoServiceInterface


class UsersServiceInterface(Protocol):
    repo: UsersRepositoryInterface

    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> SomeUser:
        raise NotImplementedError

    async def create_user(self, new_user: NewUserRequest) -> User:
        raise NotImplementedError

    async def update_mfa(self, user_id: Identifier, updates: UserMFAUpdateRequest):
        raise NotImplementedError

    async def reset_mfa(self, user_id: Identifier):
        raise NotImplementedError

    async def check_backup_code(self, user_backup_codes: str, given_backup_code: str):
        raise NotImplementedError

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        raise NotImplementedError


class UsersService(UsersServiceInterface):
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

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        return UserResponse(**user.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        return UsersResponse(
            total=total,
            items=[UserResponse(**user.dict()) for user in user_list],
        )
