from typing import Protocol

from sqlalchemy import insert, select

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.idp import people
from observer.entities.idp import IDP, NewIDP, UpdateIDP


class IDPRepositoryInterface(Protocol):
    async def create_idp(self, new_idp: NewIDP) -> IDP:
        raise NotImplementedError

    async def get_idp(self, idp_id: Identifier) -> IDP | None:
        raise NotImplementedError

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDP) -> IDP:
        raise NotImplementedError


class IDPRepository(IDPRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def create_idp(self, new_idp: NewIDP) -> IDP:
        query = insert(people).values(**new_idp.dict()).returning("*")
        result = await self.db.fetchone(query)
        return IDP(**result)

    async def get_idp(self, idp_id: Identifier) -> IDP | None:
        query = select(people).where(people.c.id == idp_id)
        if result := await self.db.fetchone(query):
            return IDP(**result)

        return None

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDP) -> IDP:
        raise NotImplementedError
