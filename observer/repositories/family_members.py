from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.idp import family_members
from observer.entities.family_members import (
    FamilyMember,
    NewFamilyMember,
    UpdateFamilyMember,
)


class IFamilyRepository(Protocol):
    async def add_member(self, new_member: NewFamilyMember) -> FamilyMember:
        raise NotImplementedError

    async def get_member(self, member_id: Identifier) -> Optional[FamilyMember]:
        raise NotImplementedError

    async def get_by_person(self, idp_id: Identifier) -> List[FamilyMember]:
        raise NotImplementedError

    async def get_by_project(self, project_id: Identifier) -> List[FamilyMember]:
        raise NotImplementedError

    async def update_member(self, member_id: Identifier, updates: UpdateFamilyMember) -> FamilyMember:
        raise NotImplementedError

    async def delete_member(self, member_id: Identifier) -> FamilyMember:
        raise NotImplementedError


class FamilyRepository(IFamilyRepository):
    def __init__(self, db: Database):
        self.db = db

    async def add_member(self, new_member: NewFamilyMember) -> FamilyMember:
        query = insert(family_members).values(**new_member.dict()).returning("*")
        result = await self.db.fetchone(query)
        return FamilyMember(**result)

    async def get_member(self, member_id: Identifier) -> Optional[FamilyMember]:
        query = select(family_members).where(family_members.c.id == member_id)
        if result := await self.db.fetchone(query):
            return FamilyMember(**result)

        return None

    async def get_by_person(self, idp_id: Identifier) -> List[FamilyMember]:
        query = select(family_members).where(family_members.c.idp_id == idp_id)
        rows = await self.db.fetchall(query)
        return [FamilyMember(**row) for row in rows]

    async def get_by_project(self, project_id: Identifier) -> List[FamilyMember]:
        query = select(family_members).where(family_members.c.project_id == project_id)
        rows = await self.db.fetchall(query)
        return [FamilyMember(**row) for row in rows]

    async def update_member(self, member_id: Identifier, updates: UpdateFamilyMember) -> FamilyMember:
        query = update(family_members).values(**updates.dict()).where(family_members.c.id == member_id).returning("*")
        result = await self.db.fetchone(query)
        return FamilyMember(**result)

    async def delete_member(self, member_id: Identifier) -> FamilyMember:
        query = delete(family_members).where(family_members.c.id == member_id).returning("*")
        result = await self.db.fetchone(query)
        return FamilyMember(**result)
