from typing import List, Optional, Protocol, Tuple

from sqlalchemy import delete, func, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.pets import pets
from observer.entities.pets import NewPet, Pet, UpdatePet
from observer.schemas.pagination import Pagination


class IPetsRepository(Protocol):
    async def create_pet(self, new_pet: NewPet) -> Pet:
        raise NotImplementedError

    async def get_pet(self, pet_id: Identifier) -> Optional[Pet]:
        raise NotImplementedError

    async def get_pets_by_project(self, project_id: Identifier, page: Pagination) -> Tuple[int, List[Pet]]:
        raise NotImplementedError

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Optional[Pet]:
        raise NotImplementedError

    async def delete_pet(self, pet_id: Identifier) -> Optional[Pet]:
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

    async def get_pets_by_project(self, project_id: Identifier, page: Pagination) -> Tuple[int, List[Pet]]:
        count_query = (
            select(
                func.count().label("count"),
            )
            .select_from(pets)
            .where(
                pets.c.project_id == project_id,
            )
        )
        query = (
            select(pets)
            .where(
                pets.c.project_id == project_id,
            )
            .offset(page.offset)
            .limit(page.limit)
        )
        rows = await self.db.fetchall(query)
        items = [Pet(**row) for row in rows]
        pets_count = await self.db.fetchone(count_query)
        return pets_count["count"], items

    async def update_pet(self, pet_id: Identifier, updates: UpdatePet) -> Optional[Pet]:
        query = (
            update(pets)
            .values(**updates.dict())
            .where(
                pets.c.id == pet_id,
            )
            .returning("*")
        )
        if row := await self.db.fetchone(query):
            return Pet(**row)

        return None

    async def delete_pet(self, pet_id: Identifier) -> Optional[Pet]:
        query = delete(pets).where(pets.c.id == pet_id).returning("*")
        if row := await self.db.fetchone(query):
            return Pet(**row)

        return None
