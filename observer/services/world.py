from datetime import datetime, timedelta, timezone
from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.base import SomeCountry, SomeState
from observer.entities.world import (
    Country,
    NewCountry,
    NewState,
    State,
    UpdateCountry,
    UpdateState,
)
from observer.repositories.world import WorldRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.world import (
    CountryResponse,
    NewCountryRequest,
    NewStateRequest,
    StateResponse,
    UpdateCountryRequest,
    UpdateStateRequest,
)


class WorldServiceInterface(Protocol):
    tag: str
    repo: WorldRepositoryInterface

    # Countries
    async def create_country(self, new_country: NewCountryRequest) -> Country:
        raise NotImplementedError

    async def get_countries(self) -> List[Country]:
        raise NotImplementedError

    async def get_country(self, country_id: Identifier) -> SomeCountry:
        raise NotImplementedError

    async def update_country(self, country_id: Identifier, updates: UpdateCountryRequest) -> SomeCountry:
        raise NotImplementedError

    async def delete_country(self, country_id: Identifier) -> SomeCountry:
        raise NotImplementedError

    @staticmethod
    async def country_to_response(country: Country) -> CountryResponse:
        raise NotImplementedError

    @staticmethod
    async def countries_to_response(country_list: List[Country]) -> List[CountryResponse]:
        raise NotImplementedError

    # States
    async def create_state(self, new_state: NewStateRequest) -> State:
        raise NotImplementedError

    async def get_states(self) -> List[State]:
        raise NotImplementedError

    async def get_state(self, state_id: Identifier) -> SomeState:
        raise NotImplementedError

    async def update_state(self, state_id: Identifier, updates: UpdateStateRequest) -> SomeState:
        raise NotImplementedError

    async def delete_state(self, state_id: Identifier) -> SomeState:
        raise NotImplementedError

    @staticmethod
    async def state_to_response(state: State) -> StateResponse:
        raise NotImplementedError

    @staticmethod
    async def states_to_response(state_list: List[State]) -> List[StateResponse]:
        raise NotImplementedError

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError


class WorldService(WorldServiceInterface):
    tag: str = "source=service:world"

    def __init__(self, places_repository: WorldRepositoryInterface):
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
        return await self.repo.update_country(country_id, UpdateCountry(**updates.dict()))

    async def delete_country(self, country_id: Identifier) -> Country:
        return await self.repo.delete_country(country_id)

    @staticmethod
    async def country_to_response(country: Country) -> CountryResponse:
        return CountryResponse(**country.dict())

    @staticmethod
    async def countries_to_response(country_list: List[Country]) -> List[CountryResponse]:
        return [CountryResponse(**country.dict()) for country in country_list]

    # States
    async def create_state(self, new_state: NewStateRequest) -> State:
        return await self.repo.create_state(NewState(**new_state.dict()))

    async def get_states(self) -> List[State]:
        return await self.repo.get_states()

    async def get_state(self, state_id: Identifier) -> SomeState:
        if state := await self.repo.get_state(state_id):
            return state

        raise NotFoundError(message="State not found")

    async def update_state(self, state_id: Identifier, updates: UpdateStateRequest) -> SomeState:
        return await self.repo.update_state(state_id, UpdateState(**updates.dict()))

    async def delete_state(self, state_id: Identifier) -> SomeState:
        return await self.repo.delete_state(state_id)

    @staticmethod
    async def state_to_response(state: State) -> StateResponse:
        return StateResponse(**state.dict())

    @staticmethod
    async def states_to_response(state_list: List[State]) -> List[StateResponse]:
        return [StateResponse(**state.dict()) for state in state_list]

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if expires_in:
            expires_at = now + expires_in

        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=expires_at,
        )
