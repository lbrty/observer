from typing import List, Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier, PlaceFilters, StateFilters
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
from observer.repositories.world import IWorldRepository
from observer.schemas.world import (
    NewCountryRequest,
    NewPlaceRequest,
    NewStateRequest,
    UpdateCountryRequest,
    UpdatePlaceRequest,
    UpdateStateRequest,
)


class IWorldService(Protocol):
    repo: IWorldRepository

    # Countries
    async def create_country(self, new_country: NewCountryRequest) -> Country:
        raise NotImplementedError

    async def get_countries(self) -> List[Country]:
        raise NotImplementedError

    async def get_country(self, country_id: Identifier) -> Country:
        raise NotImplementedError

    async def update_country(self, country_id: Identifier, updates: UpdateCountryRequest) -> Country:
        raise NotImplementedError

    async def delete_country(self, country_id: Identifier) -> Country:
        raise NotImplementedError

    # States
    async def create_state(self, new_state: NewStateRequest) -> State:
        raise NotImplementedError

    async def get_states(self, filters: Optional[StateFilters]) -> List[State]:
        raise NotImplementedError

    async def get_state(self, state_id: Identifier) -> State:
        raise NotImplementedError

    async def update_state(self, state_id: Identifier, updates: UpdateStateRequest) -> State:
        raise NotImplementedError

    async def delete_state(self, state_id: Identifier) -> State:
        raise NotImplementedError

    # Places
    async def create_place(self, new_place: NewPlaceRequest) -> Place:
        raise NotImplementedError

    async def get_places(self, filters: Optional[PlaceFilters]) -> List[Place]:
        raise NotImplementedError

    async def get_place(self, place_id: Identifier) -> Place:
        raise NotImplementedError

    async def update_place(self, place_id: Identifier, updates: UpdatePlaceRequest) -> Place:
        raise NotImplementedError

    async def delete_place(self, place_id: Identifier) -> Place:
        raise NotImplementedError


class WorldService(IWorldService):
    def __init__(self, places_repository: IWorldRepository):
        self.repo = places_repository

    # States
    async def create_country(self, new_country: NewCountryRequest) -> Country:
        return await self.repo.create_country(NewCountry(**new_country.dict()))

    async def get_countries(self) -> List[Country]:
        return await self.repo.get_countries()

    async def get_country(self, country_id: Identifier) -> Country:
        if country := await self.repo.get_country(country_id):
            return country

        raise NotFoundError(message="Country not found")

    async def update_country(self, country_id: Identifier, updates: UpdateCountryRequest) -> Country:
        if country := await self.repo.update_country(country_id, UpdateCountry(**updates.dict())):
            return country

        raise NotFoundError(message="Country not found")

    async def delete_country(self, country_id: Identifier) -> Country:
        if country := await self.repo.delete_country(country_id):
            return country

        raise NotFoundError(message="Country not found")

    # States
    async def create_state(self, new_state: NewStateRequest) -> State:
        return await self.repo.create_state(NewState(**new_state.dict()))

    async def get_states(self, filters: Optional[StateFilters]) -> List[State]:
        return await self.repo.get_states(filters)

    async def get_state(self, state_id: Identifier) -> State:
        if state := await self.repo.get_state(state_id):
            return state

        raise NotFoundError(message="State not found")

    async def update_state(self, state_id: Identifier, updates: UpdateStateRequest) -> State:
        if state := await self.repo.update_state(state_id, UpdateState(**updates.dict())):
            return state

        raise NotFoundError(message="State not found")

    async def delete_state(self, state_id: Identifier) -> State:
        if state := await self.repo.delete_state(state_id):
            return state

        raise NotFoundError(message="State not found")

    # Places
    async def create_place(self, new_place: NewPlaceRequest) -> Place:
        return await self.repo.create_place(NewPlace(**new_place.dict()))

    async def get_places(self, filters: Optional[PlaceFilters]) -> List[Place]:
        return await self.repo.get_places(filters)

    async def get_place(self, place_id: Identifier) -> Place:
        if place := await self.repo.get_place(place_id):
            return place

        raise NotFoundError(message="Place not found")

    async def update_place(self, place_id: Identifier, updates: UpdatePlaceRequest) -> Place:
        if place := await self.repo.update_place(place_id, UpdatePlace(**updates.dict())):
            return place

        raise NotFoundError(message="Place not found")

    async def delete_place(self, place_id: Identifier) -> Place:
        if place := await self.repo.delete_place(place_id):
            return place

        raise NotFoundError(message="Place not found")
