from typing import Protocol

from sqlalchemy import insert

from observer.db import Database
from observer.db.tables.idp import people
from observer.entities.idp import IDP, NewIDP


class IDPRepositoryInterface(Protocol):
    async def create_idp(self, new_idp: NewIDP) -> IDP:
        raise NotImplementedError


class IDPRepository(IDPRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def create_idp(self, new_idp: NewIDP) -> IDP:
        query = insert(people).values(**new_idp.dict()).returning("*")
        result = await self.db.fetchone(query)
        return IDP(**result)
