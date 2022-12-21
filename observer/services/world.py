from datetime import datetime, timedelta, timezone
from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.world import Country, NewCountry, UpdateCountry
from observer.repositories.world import WorldRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.world import (
    CountryResponse,
    NewCountryRequest,
    UpdateCountryRequest,
)


class WorldServiceInterface(Protocol):
    tag: str
    repo: WorldRepositoryInterface

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

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError

    @staticmethod
    async def country_to_response(country: Country) -> CountryResponse:
        raise NotImplementedError

    @staticmethod
    async def countries_to_response(country_list: List[Country]) -> List[CountryResponse]:
        raise NotImplementedError


class WorldService(WorldServiceInterface):
    tag: str = "source=service:places"

    def __init__(self, places_repository: WorldRepositoryInterface):
        self.repo = places_repository

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

    @staticmethod
    async def country_to_response(country: Country) -> CountryResponse:
        return CountryResponse(**country.dict())

    @staticmethod
    async def countries_to_response(country_list: List[Country]) -> List[CountryResponse]:
        return [CountryResponse(**country.dict()) for country in country_list]
