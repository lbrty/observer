from typing import List, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.world import countries
from observer.entities.world import Country, NewCountry, UpdateCountry


class WorldRepositoryInterface(Protocol):
    async def create_country(self, new_country: NewCountry) -> Country:
        raise NotImplementedError

    async def get_countries(self) -> List[Country]:
        raise NotImplementedError

    async def get_country(self, country_id: Identifier) -> Country:
        raise NotImplementedError

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> Country:
        raise NotImplementedError

    async def delete_country(self, country_id: Identifier) -> Country:
        raise NotImplementedError


class WorldRepository(WorldRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def create_country(self, new_country: NewCountry) -> Country:
        query = insert(countries).values(**new_country.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Country(**result)

    async def get_countries(self) -> List[Country]:
        query = select(countries)
        rows = await self.db.fetchall(query)
        return [Country(**row) for row in rows]

    async def get_country(self, country_id: Identifier) -> Country:
        query = select(countries).where(countries.c.id == country_id)
        row = await self.db.fetchone(query)
        return Country(**row)

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> Country:
        query = update(countries).values(**updates.dict()).where(countries.c.id == country_id).returning("*")
        row = await self.db.fetchone(query)
        return Country(**row)

    async def delete_country(self, country_id: Identifier) -> Country:
        query = delete(countries).where(countries.c.id == country_id).returning("*")
        row = await self.db.fetchone(query)
        return Country(**row)
