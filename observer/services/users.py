from datetime import datetime, timedelta, timezone
from typing import List, Optional, Protocol, Tuple

import shortuuid

from observer.api.exceptions import (
    ConfirmationCodeExpiredError,
    InviteExpiredError,
    NotFoundError,
    TOTPInvalidBackupCodeError,
)
from observer.common import bcrypt
from observer.common.types import Identifier, Pagination, UserFilters
from observer.entities.users import (
    Confirmation,
    Invite,
    NewUser,
    PasswordReset,
    User,
    UserUpdate,
)
from observer.repositories.users import IUsersRepository
from observer.schemas.users import NewUserRequest, UserMFAUpdateRequest
from observer.services.crypto import ICryptoService
from observer.settings import settings


class IUsersService(Protocol):
    repo: IUsersRepository

    async def get_by_id(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> Optional[User]:
        raise NotImplementedError

    async def create_user(self, new_user: NewUserRequest, is_active: bool = True) -> User:
        raise NotImplementedError

    async def delete_user(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        raise NotImplementedError

    async def deactivate_user(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def update_mfa(self, user_id: Identifier, updates: UserMFAUpdateRequest):
        raise NotImplementedError

    async def reset_mfa(self, user_id: Identifier):
        raise NotImplementedError

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> User:
        raise NotImplementedError

    async def filter_users(self, filters: UserFilters, pages: Pagination) -> Tuple[int, List[User]]:
        raise NotImplementedError

    async def check_backup_code(self, user_backup_codes: str, given_backup_code: str):
        raise NotImplementedError

    async def confirm_user(self, user_id: Optional[Identifier], code: str) -> User:
        raise NotImplementedError

    async def just_confirm_user(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def create_confirmation(self, user_id: Identifier) -> Confirmation:
        raise NotImplementedError

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        raise NotImplementedError

    async def reset_password(self, user_id: Identifier) -> PasswordReset:
        raise NotImplementedError

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        raise NotImplementedError

    async def create_invite(self, user_id: Identifier) -> Invite:
        raise NotImplementedError

    async def get_invite(self, code: str, validate: bool = True) -> Invite:
        raise NotImplementedError

    async def get_invites(self, page: Pagination) -> Tuple[int, List[Invite]]:
        raise NotImplementedError

    async def delete_invite(self, code: str) -> Invite:
        raise NotImplementedError


class UsersService(IUsersService):
    def __init__(self, users_repository: IUsersRepository, crypto_service: ICryptoService):
        self.repo = users_repository
        self.crypto_service = crypto_service

    async def get_by_id(self, user_id: Identifier) -> User:
        if user := await self.repo.get_by_id(user_id):
            return user

        raise NotFoundError(message="User not found")

    async def get_by_email(self, email: str) -> Optional[User]:
        return await self.repo.get_by_email(email)

    async def create_user(self, new_user: NewUserRequest, is_active: bool = True) -> User:
        password_hash = bcrypt.hash_password(new_user.password.get_secret_value())
        user = NewUser(
            email=new_user.email,
            full_name=new_user.full_name,
            password_hash=password_hash,
            role=new_user.role,
            is_active=is_active,
            is_confirmed=False,
        )
        return await self.repo.create_user(user)

    async def delete_user(self, user_id: Identifier) -> User:
        if user := await self.repo.delete_user(user_id):
            return user

        raise NotFoundError(message="User not found")

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        if user := await self.repo.update_password(user_id, new_password_hash):
            return user

        raise NotFoundError(message="User not found")

    async def deactivate_user(self, user_id: Identifier) -> User:
        if user := await self.repo.update_user(user_id, UserUpdate(is_active=False)):
            return user

        raise NotFoundError(message="User not found")

    async def update_mfa(self, user_id: Identifier, updates: UserMFAUpdateRequest):
        user_update = UserUpdate(
            mfa_enabled=updates.mfa_enabled,
            mfa_encrypted_secret=updates.mfa_encrypted_secret,
            mfa_encrypted_backup_codes=updates.mfa_encrypted_backup_codes,
        )
        if user := await self.repo.update_user(user_id, user_update):
            return user

        raise NotFoundError(message="User not found")

    async def reset_mfa(self, user_id: Identifier):
        if user := await self.repo.reset_mfa(user_id):
            return user

        raise NotFoundError(message="User not found")

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> User:
        if user := await self.repo.update_user(user_id, updates):
            return user

        raise NotFoundError(message="User not found")

    async def filter_users(self, filters: UserFilters, pages: Pagination) -> Tuple[int, List[User]]:
        return await self.repo.filter_users(filters, pages)

    async def check_backup_code(self, user_backup_codes: str, given_backup_code: str):
        keys_hash, encrypted_backup_codes = user_backup_codes.split(":", maxsplit=1)
        decrypted_backup_codes = await self.crypto_service.decrypt(keys_hash, encrypted_backup_codes.encode())
        if given_backup_code not in decrypted_backup_codes.decode().split(","):
            raise TOTPInvalidBackupCodeError(message="Invalid backup code")

    async def confirm_user(self, user_id: Optional[Identifier], code: str) -> User:
        confirmation = await self.get_confirmation(code)

        # If user is authenticated then we need to check
        # if confirmation code belongs to this user
        # and if not so then we need to return not found error.
        if not confirmation or (user_id is not None and confirmation.user_id != user_id):
            raise NotFoundError(message="Confirmation code not found")

        if confirmation.expires_at < datetime.now(tz=timezone.utc):
            raise ConfirmationCodeExpiredError(message="Confirmation code has already expired")

        # For the case if confirmation code has not expired, and yet we still
        # get requests we just need to return original confirmation instance.
        user = await self.get_by_id(confirmation.user_id)
        if not user:
            raise NotFoundError(message="Unknown user")

        if user.is_confirmed:
            return user

        await self.repo.confirm_user(confirmation.user_id)
        return user

    async def just_confirm_user(self, user_id: Identifier) -> User:
        if user := await self.repo.confirm_user(user_id):
            return user

        raise NotFoundError(message="User not found")

    async def create_confirmation(self, user_id: Identifier) -> Confirmation:
        now = datetime.now(tz=timezone.utc)
        confirmation_delta = timedelta(minutes=settings.confirmation_expiration_minutes)
        return await self.repo.create_confirmation(user_id, shortuuid.uuid(), now + confirmation_delta)

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        return await self.repo.get_confirmation(code)

    async def create_invite(self, user_id: Identifier) -> Invite:
        now = datetime.now(tz=timezone.utc)
        delta = timedelta(minutes=settings.invite_expiration_minutes)
        return await self.repo.create_invite(user_id, shortuuid.uuid(), now + delta)

    async def get_invite(self, code: str, validate: bool = True) -> Invite:
        invite = await self.repo.get_invite(code)

        # If user is authenticated then we need to check
        # if invite code belongs to this user
        # and if not so then we need to return not found error.
        if not invite:
            raise NotFoundError(message="Invite not found")

        if not validate:
            return invite

        if invite.expires_at < datetime.now(tz=timezone.utc):
            raise InviteExpiredError(message="Invite has already expired")

        # For the case if confirmation code has not expired, and yet we still
        # get requests we just need to return original confirmation instance.
        user = await self.get_by_id(invite.user_id)
        if not user:
            raise NotFoundError(message="Unknown user")

        return invite

    async def get_invites(self, page: Pagination) -> Tuple[int, List[Invite]]:
        return await self.repo.get_invites(page.offset, page.limit)

    async def delete_invite(self, code: str) -> Invite:
        if invite := await self.repo.delete_invite(code):
            return invite

        raise NotFoundError(message="Invite not found")

    async def reset_password(self, user_id: Identifier) -> PasswordReset:
        return await self.repo.create_password_reset_code(user_id, shortuuid.uuid())

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        return await self.repo.get_password_reset(code)
