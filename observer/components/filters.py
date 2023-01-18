from typing import Optional

from fastapi import Query

from observer.common.types import Identifier, PlaceFilters, PlaceType, StateFilters


async def state_filters(
    name: Optional[str] = Query(None, description="Name of state"),
    code: Optional[str] = Query(None, description="Code of state"),
    country_id: Optional[Identifier] = Query(None, description="ID of country"),
) -> StateFilters:
    return StateFilters(
        name=name,
        code=code,
        country_id=country_id,
    )


async def place_filters(
    name: Optional[str] = Query(None, description="Name of place"),
    code: Optional[str] = Query(None, description="Code of place"),
    place_type: PlaceType | None = Query(None, description="Type of place"),
    state_id: Optional[Identifier] = Query(None, description="ID of state"),
    country_id: Optional[Identifier] = Query(None, description="ID of country"),
) -> PlaceFilters:
    return PlaceFilters(
        name=name,
        code=code,
        place_type=place_type,
        state_id=state_id,
        country_id=country_id,
    )
