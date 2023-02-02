from typing import List, Optional, Protocol

from sqlalchemy import and_, delete, insert, select, update

from observer.common.types import Identifier, PlaceFilters, StateFilters
from observer.db import Database
from observer.db.tables.world import countries, places, states
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


class IWorldRepository(Protocol):
    async def create_country(self, new_country: NewCountry) -> Country:
        raise NotImplementedError

    # Countries
    async def get_countries(self) -> List[Country]:
        raise NotImplementedError

    async def get_country(self, country_id: Identifier) -> Optional[Country]:
        raise NotImplementedError

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> Optional[Country]:
        raise NotImplementedError

    async def delete_country(self, country_id: Identifier) -> Optional[Country]:
        raise NotImplementedError

    # States
    async def create_state(self, new_state: NewState) -> State:
        raise NotImplementedError

    async def get_states(self, filters: Optional[StateFilters]) -> List[State]:
        raise NotImplementedError

    async def get_state(self, country_id: Identifier) -> Optional[State]:
        raise NotImplementedError

    async def update_state(self, country_id: Identifier, updates: UpdateState) -> Optional[State]:
        raise NotImplementedError

    async def delete_state(self, country_id: Identifier) -> Optional[State]:
        raise NotImplementedError

    # Places
    async def create_place(self, new_place: NewPlace) -> Place:
        raise NotImplementedError

    async def get_places(self, filters: Optional[PlaceFilters]) -> List[Place]:
        raise NotImplementedError

    async def get_place(self, place_id: Identifier) -> Optional[Place]:
        raise NotImplementedError

    async def update_place(self, place_id: Identifier, updates: UpdatePlace) -> Optional[Place]:
        raise NotImplementedError

    async def delete_place(self, place_id: Identifier) -> Optional[Place]:
        raise NotImplementedError


class WorldRepository(IWorldRepository):
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

    async def get_country(self, country_id: Identifier) -> Optional[Country]:
        query = select(countries).where(countries.c.id == country_id)
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    async def update_country(self, country_id: Identifier, updates: UpdateCountry) -> Optional[Country]:
        query = update(countries).values(**updates.dict()).where(countries.c.id == country_id).returning("*")
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    async def delete_country(self, country_id: Identifier) -> Optional[Country]:
        query = delete(countries).where(countries.c.id == country_id).returning("*")
        if row := await self.db.fetchone(query):
            return Country(**row)

        return None

    # States
    async def create_state(self, new_state: NewState) -> State:
        query = insert(states).values(**new_state.dict()).returning("*")
        result = await self.db.fetchone(query)
        return State(**result)

    async def get_states(self, filters: Optional[StateFilters]) -> List[State]:
        conditions = []
        if filters.name:
            conditions.append(states.c.name.ilike(f"%{filters.name}%"))

        if filters.code:
            conditions.append(states.c.code.ilike(f"%{filters.code}%"))

        if filters.country_id:
            conditions.append(states.c.country_id == filters.country_id)

        query = select(states)
        if len(conditions) > 0:
            query = query.where(and_(*conditions))

        rows = await self.db.fetchall(query)
        return [State(**row) for row in rows]

    async def get_state(self, state_id: Identifier) -> Optional[State]:
        query = select(states).where(states.c.id == state_id)
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    async def update_state(self, state_id: Identifier, updates: UpdateState) -> Optional[State]:
        query = update(states).values(**updates.dict()).where(states.c.id == state_id).returning("*")
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    async def delete_state(self, state_id: Identifier) -> Optional[State]:
        query = delete(states).where(states.c.id == state_id).returning("*")
        if row := await self.db.fetchone(query):
            return State(**row)

        return None

    # Places
    async def create_place(self, new_place: NewPlace) -> Place:
        query = insert(places).values(**new_place.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Place(**result)

    async def get_places(self, filters: Optional[PlaceFilters]) -> List[Place]:
        conditions = []
        if filters.name:
            conditions.append(places.c.name.ilike(f"%{filters.name}%"))

        if filters.code:
            conditions.append(places.c.code.ilike(f"%{filters.code}%"))

        if filters.place_type:
            conditions.append(places.c.place_type == filters.place_type)

        if filters.state_id:
            conditions.append(places.c.state_id == filters.state_id)

        if filters.country_id:
            conditions.append(places.c.country_id == filters.country_id)

        query = select(places)
        if len(conditions) > 0:
            query = query.where(and_(*conditions))

        rows = await self.db.fetchall(query)
        return [Place(**row) for row in rows]

    async def get_place(self, place_id: Identifier) -> Optional[Place]:
        query = select(places).where(places.c.id == place_id)
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None

    async def update_place(self, place_id: Identifier, updates: UpdatePlace) -> Optional[Place]:
        query = update(places).values(**updates.dict()).where(places.c.id == place_id).returning("*")
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None

    async def delete_place(self, place_id: Identifier) -> Optional[Place]:
        query = delete(places).where(places.c.id == place_id).returning("*")
        if row := await self.db.fetchone(query):
            return Place(**row)

        return None
