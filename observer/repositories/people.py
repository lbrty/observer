from typing import Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import EncryptedFieldValue, Identifier
from observer.db import Database
from observer.db.tables.people import people
from observer.entities.people import NewPerson, Person, UpdatePerson


class IPeopleRepository(Protocol):
    async def create_person(self, new_person: NewPerson) -> Person:
        raise NotImplementedError

    async def get_person(self, person_id: Identifier) -> Optional[Person]:
        raise NotImplementedError

    async def update_person(self, person_id: Identifier, updates: UpdatePerson) -> Optional[Person]:
        raise NotImplementedError

    async def delete_person(self, person_id: Identifier) -> Optional[Person]:
        raise NotImplementedError


class PeopleRepository(IPeopleRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_person(self, new_person: NewPerson) -> Person:
        query = insert(people).values(**new_person.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Person(**result)

    async def get_person(self, person_id: Identifier) -> Optional[Person]:
        query = select(people).where(people.c.id == person_id)
        if result := await self.db.fetchone(query):
            return Person(**result)

        return None

    async def update_person(self, person_id: Identifier, updates: UpdatePerson) -> Optional[Person]:
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

        query = update(people).values(values).where(people.c.id == person_id).returning("*")
        if result := await self.db.fetchone(query):
            return Person(**result)

        return None

    async def delete_person(self, person_id: Identifier) -> Optional[Person]:
        query = delete(people).where(people.c.id == person_id).returning("*")
        if result := await self.db.fetchone(query):
            return Person(**result)

        return None
