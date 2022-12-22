from typing import List, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.world import countries, places, states
from observer.entities.base import SomeCountry, SomePlace, SomeState
from observer.entities.world import (
    Country,
    NewCountry,
    NewPlace,
    NewState,
    Place,
    State,
    UpdateCountry,
    UpdatePlace,
    UpdateState,
)


class WorldRepositoryInterface(Protocol):
    async def create_country(self, new_country: NewCountry) -> Country:
        raise NotImplementedError

    # Countries
    async def get_countries(self) -> List[Country]:
        raise NotImplementedError

    async def get_country(self, country_id: Identifier) -> SomeCountry:
        raise NotImplementedError

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> SomeCountry:
        raise NotImplementedError

    async def delete_country(self, country_id: Identifier) -> SomeCountry:
        raise NotImplementedError

    # States
    async def create_state(self, new_state: NewState) -> State:
        raise NotImplementedError

    async def get_states(self) -> List[State]:
        raise NotImplementedError

    async def get_state(self, country_id: Identifier) -> SomeState:
        raise NotImplementedError

    async def update_state(self, country_id: Identifier, updates: UpdateState) -> SomeState:
        raise NotImplementedError

    async def delete_state(self, country_id: Identifier) -> SomeState:
        raise NotImplementedError

    # Places
    async def create_place(self, new_place: NewPlace) -> Place:
        raise NotImplementedError

    async def get_places(self) -> List[Place]:
        raise NotImplementedError

    async def get_place(self, place_id: Identifier) -> SomePlace:
        raise NotImplementedError

    async def update_place(self, place_id: Identifier, updates: UpdatePlace) -> SomePlace:
        raise NotImplementedError

    async def delete_place(self, place_id: Identifier) -> SomePlace:
        raise NotImplementedError


class WorldRepository(WorldRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    # Countries
    async def create_country(self, new_country: NewCountry) -> Country:
        query = insert(countries).values(**new_country.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Country(**result)

    async def get_countries(self) -> List[Country]:
        query = select(countries)
        rows = await self.db.fetchall(query)
        return [Country(**row) for row in rows]

    async def get_country(self, country_id: Identifier) -> SomeCountry:
        query = select(countries).where(countries.c.id == country_id)
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> SomeCountry:
        query = update(countries).values(**updates.dict()).where(countries.c.id == country_id).returning("*")
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    async def delete_country(self, country_id: Identifier) -> SomeCountry:
        query = delete(countries).where(countries.c.id == country_id).returning("*")
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    # States
    async def create_state(self, new_state: NewState) -> State:
        query = insert(states).values(**new_state.dict()).returning("*")
        result = await self.db.fetchone(query)
        return State(**result)

    async def get_states(self) -> List[State]:
        query = select(states)
        rows = await self.db.fetchall(query)
        return [State(**row) for row in rows]

    async def get_state(self, state_id: Identifier) -> SomeState:
        query = select(states).where(states.c.id == state_id)
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    async def update_state(self, state_id: Identifier, updates: UpdateState) -> SomeState:
        query = update(states).values(**updates.dict()).where(states.c.id == state_id).returning("*")
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    async def delete_state(self, state_id: Identifier) -> SomeState:
        query = delete(states).where(states.c.id == state_id).returning("*")
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    # Places
    async def create_place(self, new_place: NewPlace) -> Place:
        query = insert(places).values(**new_place.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Place(**result)

    async def get_places(self) -> List[Place]:
        query = select(places)
        rows = await self.db.fetchall(query)
        return [Place(**row) for row in rows]

    async def get_place(self, place_id: Identifier) -> SomePlace:
        query = select(places).where(places.c.id == place_id)
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None

    async def update_place(self, place_id: Identifier, updates: UpdatePlace) -> SomePlace:
        query = update(places).values(**updates.dict()).where(places.c.id == place_id).returning("*")
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None

    async def delete_place(self, place_id: Identifier) -> SomePlace:
        query = delete(places).where(places.c.id == place_id).returning("*")
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None