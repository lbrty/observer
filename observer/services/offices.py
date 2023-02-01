from typing import List, Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.offices import Office
from observer.repositories.offices import IOfficesRepository


class IOfficesService(Protocol):
    async def create_office(self, name: str) -> Office:
        raise NotImplementedError

    async def get_office(self, office_id: Identifier) -> Optional[Office]:
        raise NotImplementedError

    async def get_offices(self, name: Optional[str], offset: int, limit: int) -> List[Office]:
        raise NotImplementedError

    async def update_office(self, office_id: Identifier, new_name: str) -> Office:
        raise NotImplementedError

    async def delete_office(self, office_id: Identifier) -> Office:
        raise NotImplementedError


class OfficesService(IOfficesService):
    def __init__(self, repo: IOfficesRepository):
        self.repo = repo

    async def create_office(self, name: str) -> Office:
        return await self.repo.create_office(name)

    async def get_office(self, office_id: Identifier) -> Optional[Office]:
        if office := await self.repo.get_office(office_id):
            return office

        raise NotFoundError(message="Office not found")

    async def get_offices(self, name: Optional[str], offset: int, limit: int) -> List[Office]:
        return await self.repo.get_offices(name, offset, limit)

    async def update_office(self, office_id: Identifier, new_name: str) -> Office:
        return await self.repo.update_office(office_id, new_name)

    async def delete_office(self, office_id: Identifier) -> Office:
        return await self.repo.delete_office(office_id)
