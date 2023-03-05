from datetime import datetime, timezone
from typing import List, Optional, Protocol, Tuple

from sqlalchemy import and_, delete, desc, func, insert, select, update

from observer.common.types import Identifier, Pagination, UserFilters
from observer.db import Database
from observer.db.tables.users import confirmations, invites, password_resets, users
from observer.entities.users import (
    Confirmation,
    Invite,
    NewUser,
    PasswordReset,
    User,
    UserUpdate,
)


class IUsersRepository(Protocol):
    async def get_by_id(self, user_id: Identifier) -> Optional[User]:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> Optional[User]:
        raise NotImplementedError

    async def create_user(self, new_user: NewUser) -> User:
        raise NotImplementedError

    async def delete_user(self, user_id: Identifier) -> Optional[User]:
        raise NotImplementedError

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> Optional[User]:
        raise NotImplementedError

    async def filter_users(self, filters: UserFilters, pages: Pagination) -> Tuple[int, List[User]]:
        raise NotImplementedError

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> Optional[User]:
        raise NotImplementedError

    async def reset_mfa(self, user_id: Identifier) -> Optional[User]:
        raise NotImplementedError

    async def create_password_reset_code(self, user_id: Identifier, code: str) -> PasswordReset:
        raise NotImplementedError

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        raise NotImplementedError

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        raise NotImplementedError

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        raise NotImplementedError

    async def confirm_user(self, user_id: Identifier) -> Optional[User]:
        raise NotImplementedError

    async def create_invite(self, user_id: Identifier, code: str, expires_at: datetime) -> Invite:
        raise NotImplementedError

    async def get_invite(self, code: str) -> Optional[Invite]:
        raise NotImplementedError

    async def get_invites(self, offset: int, limit: int) -> Tuple[int, List[Invite]]:
        raise NotImplementedError

    async def delete_invite(self, code: str) -> Optional[Invite]:
        raise NotImplementedError


class UsersRepository(IUsersRepository):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, user_id: Identifier) -> Optional[User]:
        query = select(users).where(users.c.id == user_id)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def get_by_email(self, email: str) -> Optional[User]:
        query = select(users).where(users.c.email == email)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def create_user(self, new_user: NewUser) -> User:
        query = insert(users).values(**new_user.dict()).returning("*")
        result = await self.db.fetchone(query)
        return User(**result)

    async def delete_user(self, user_id: Identifier) -> Optional[User]:
        query = delete(users).where(users.c.id == user_id).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        values = dict(
            code=code,
            user_id=user_id,
            expires_at=expires_at,
        )
        query = insert(confirmations).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return Confirmation(**result)

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        query = select(confirmations).where(confirmations.c.code == code)
        if result := await self.db.fetchone(query):
            return Confirmation(**result)

        return None

    async def create_invite(self, user_id: Identifier, code: str, expires_at: datetime) -> Invite:
        values = dict(
            code=code,
            user_id=user_id,
            expires_at=expires_at,
        )
        query = insert(invites).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return Invite(**result)

    async def get_invite(self, code: str) -> Optional[Invite]:
        query = select(invites).where(invites.c.code == code)
        if result := await self.db.fetchone(query):
            return Invite(**result)

        return None

    async def get_invites(self, offset: int, limit: int) -> Tuple[int, List[Invite]]:
        count_query = select(func.count().label("count")).select_from(invites)
        query = select(invites).offset(offset).limit(limit).order_by(desc(invites.c.expires_at))
        rows = await self.db.fetchall(query)
        items = [Invite(**row) for row in rows]
        invites_count = await self.db.fetchone(count_query)
        return invites_count["count"], items

    async def delete_invite(self, code: str) -> Optional[Invite]:
        query = delete(invites).where(invites.c.code == code).returning("*")
        if result := await self.db.fetchone(query):
            return Invite(**result)

        return None

    async def confirm_user(self, user_id: Identifier) -> Optional[User]:
        query = (
            update(users)
            .values(dict(is_confirmed=True))
            .where(
                users.c.id == user_id,
            )
            .returning("*")
        )
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> Optional[User]:
        update_values = {}
        if updates.email:
            update_values["email"] = updates.email

        if updates.full_name:
            update_values["full_name"] = updates.full_name  # type:ignore

        if updates.role:
            update_values["role"] = updates.role.value  # type:ignore

        if updates.is_active:
            update_values["is_active"] = updates.is_active  # type:ignore

        if updates.mfa_enabled is not None and updates.mfa_encrypted_secret and updates.mfa_encrypted_backup_codes:
            update_values["mfa_enabled"] = updates.mfa_enabled  # type:ignore
            update_values["mfa_encrypted_secret"] = updates.mfa_encrypted_secret  # type:ignore
            update_values["mfa_encrypted_backup_codes"] = updates.mfa_encrypted_backup_codes  # type:ignore

        query = update(users).values(update_values).where(users.c.id == user_id).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def filter_users(self, filters: Optional[UserFilters], pages: Pagination) -> Tuple[int, List[User]]:
        conditions = []
        if filters:
            if filters.email:
                conditions.append(users.c.email.ilike(f"%{filters.email}%"))

            if filters.full_name:
                conditions.append(users.c.full_name.ilike(f"%{filters.full_name}%"))

            if filters.role:
                conditions.append(users.c.role == filters.role)

            if filters.office_id:
                conditions.append(users.c.office_id == filters.office_id)

            if filters.is_active is not None:
                conditions.append(users.c.is_active == filters.is_active)

        query = select(users)
        if len(conditions) > 0:
            query = query.where(and_(*conditions))
        query = query.offset(pages.offset).limit(pages.limit)

        count_query = (
            select(
                func.count().label("count"),
            )
            .select_from(users)
            .where(and_(*conditions))
        )
        rows = await self.db.fetchall(query)
        items = [User(**row) for row in rows]
        users_count = await self.db.fetchone(count_query)
        return users_count["count"], items

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> Optional[User]:
        update_values = dict(password_hash=new_password_hash)
        query = update(users).values(update_values).where(users.c.id == user_id).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def reset_mfa(self, user_id: Identifier) -> Optional[User]:
        update_values = dict(
            mfa_enabled=False,
            mfa_encrypted_secret=None,
            mfa_encrypted_backup_codes=None,
        )

        query = update(users).values(update_values).where(users.c.id == user_id).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def create_password_reset_code(self, user_id: Identifier, code: str) -> PasswordReset:
        values = dict(
            code=code,
            user_id=user_id,
            created_at=datetime.now(tz=timezone.utc),
        )
        query = insert(password_resets).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return PasswordReset(**result)

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        query = select(password_resets).where(password_resets.c.code == code)
        if result := await self.db.fetchone(query):
            return PasswordReset(**result)

        return None
