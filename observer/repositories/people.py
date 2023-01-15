from typing import Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import EncryptedFieldValue, Identifier
from observer.db import Database
from observer.db.tables.people import people
from observer.entities.people import IDP, NewIDP, UpdateIDP


class IIDPRepository(Protocol):
    async def create_idp(self, new_idp: NewIDP) -> IDP:
        raise NotImplementedError

    async def get_idp(self, idp_id: Identifier) -> Optional[IDP]:
        raise NotImplementedError

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDP) -> IDP:
        raise NotImplementedError

    async def delete_idp(self, idp_id: Identifier) -> Optional[IDP]:
        raise NotImplementedError


class IDPRepository(IIDPRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_idp(self, new_idp: NewIDP) -> IDP:
        query = insert(people).values(**new_idp.dict()).returning("*")
        result = await self.db.fetchone(query)
        return IDP(**result)

    async def get_idp(self, idp_id: Identifier) -> Optional[IDP]:
        query = select(people).where(people.c.id == idp_id)
        if result := await self.db.fetchone(query):
            return IDP(**result)

        return None

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDP) -> IDP:
        values = updates.dict()
        if updates.email == EncryptedFieldValue:
            del values["email"]
        else:
            values["email"] = updates.email

        if updates.phone_number == EncryptedFieldValue:
            del values["phone_number"]
        else:
            values["phone_number"] = updates.phone_number

        if updates.phone_number_additional == EncryptedFieldValue:
            del values["phone_number_additional"]
        else:
            values["phone_number_additional"] = updates.phone_number_additional

        query = update(people).values(values).where(people.c.id == idp_id).returning("*")
        result = await self.db.fetchone(query)
        return IDP(**result)

    async def delete_idp(self, idp_id: Identifier) -> Optional[IDP]:
        query = delete(people).where(people.c.id == idp_id).returning("*")
        if result := await self.db.fetchone(query):
            return IDP(**result)

        return None