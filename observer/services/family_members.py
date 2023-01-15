from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.family_members import (
    FamilyMember,
    NewFamilyMember,
    UpdateFamilyMember,
)
from observer.repositories.family_members import IFamilyRepository
from observer.schemas.family_members import (
    NewFamilyMemberRequest,
    UpdateFamilyMemberRequest,
)


class IFamilyService(Protocol):
    repo: IFamilyRepository

    async def add_member(self, new_member: NewFamilyMemberRequest) -> FamilyMember:
        raise NotImplementedError

    async def get_member(self, member_id: Identifier) -> FamilyMember:
        raise NotImplementedError

    async def get_by_person(self, idp_id: Identifier) -> List[FamilyMember]:
        raise NotImplementedError

    async def get_by_project(self, project_id: Identifier) -> List[FamilyMember]:
        raise NotImplementedError

    async def update_member(self, member_id: Identifier, updates: UpdateFamilyMemberRequest) -> FamilyMember:
        raise NotImplementedError

    async def delete_member(self, member_id: Identifier) -> FamilyMember:
        raise NotImplementedError


class FamilyService(IFamilyService):
    def __init__(self, repo: IFamilyRepository):
        self.repo = repo

    async def add_member(self, new_member: NewFamilyMemberRequest) -> FamilyMember:
        return await self.repo.add_member(NewFamilyMember(**new_member.dict()))

    async def get_member(self, member_id: Identifier) -> FamilyMember:
        if member := await self.repo.get_member(member_id):
            return member

        raise NotFoundError(message="Family member not found")

    async def get_by_person(self, idp_id: Identifier) -> List[FamilyMember]:
        return await self.repo.get_by_person(idp_id)

    async def get_by_project(self, project_id: Identifier) -> List[FamilyMember]:
        return await self.repo.get_by_project(project_id)

    async def update_member(self, member_id: Identifier, updates: UpdateFamilyMemberRequest) -> FamilyMember:
        return await self.repo.update_member(member_id, UpdateFamilyMember(**updates.dict()))

    async def delete_member(self, member_id: Identifier) -> FamilyMember:
        return await self.repo.delete_member(member_id)
