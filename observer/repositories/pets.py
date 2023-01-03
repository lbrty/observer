from typing import Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.pets import pets
from observer.entities.pets import NewPet, Pet, UpdatePet


class IPetsRepository(Protocol):
    async def create_pet(self, new_pet: NewPet) -> Pet:
        raise NotImplementedError

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Pet:
        raise NotImplementedError

    async def delete_pet(self, pet_id: Identifier) -> Pet:
        raise NotImplementedError


class PetsRepository(IPetsRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_pet(self, new_pet: NewPet) -> Pet:
        query = insert(pets).values(**new_pet.dict()).returning("*")
        row = await self.db.fetchone(query)
        return Pet(**row)

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        query = select(pets).where(pets.c.id == pet_id)
        if row := await self.db.fetchone(query):
            return Pet(**row)

        return None

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Pet:
        query = update(pets).values(**updates.dict()).where(pets.c.id == pet_id)
        row = await self.db.fetchone(query)
        return Pet(**row)

    async def delete_pet(self, pet_id: Identifier) -> Pet:
        query = delete(pets).where(pets.c.id == pet_id).returning("*")
        row = await self.db.fetchone(query)
        return Pet(**row)
