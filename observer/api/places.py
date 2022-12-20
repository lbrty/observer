from typing import List

from fastapi import APIRouter, Depends, Response
from starlette import status

from observer.common.types import Identifier, Role
from observer.components.auth import RequiresRoles, current_user
from observer.entities.base import SomeUser
from observer.schemas.places import (
    CountryResponse,
    NewCountryRequest,
    NewPlaceRequest,
    NewStateRequest,
    PlaceResponse,
    StateResponse,
    UpdateCountryRequest,
    UpdatePlaceRequest,
    UpdateStateRequest,
)

router = APIRouter(prefix="/places")


# Countries
@router.post(
    "/countries",
    response_model=CountryResponse,
    status_code=status.HTTP_201_CREATED,
)
async def create_country(
    new_country: NewCountryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> CountryResponse:
    pass


@router.get(
    "/countries",
    response_model=List[CountryResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_countries() -> List[CountryResponse]:
    pass


@router.get(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_country(country_id: Identifier) -> CountryResponse:
    pass


@router.put(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
)
async def update_country(
    country_id: Identifier,
    updates: UpdateCountryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> CountryResponse:
    pass


@router.delete(
    "/countries/{country_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_country(
    country_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> Response:
    pass


# States
@router.post(
    "/states",
    status_code=status.HTTP_201_CREATED,
    response_model=StateResponse,
)
async def create_state(
    new_state: NewStateRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> StateResponse:
    pass


@router.get(
    "/states",
    response_model=List[StateResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_states() -> List[StateResponse]:
    pass


@router.get(
    "/states/{state_id}",
    response_model=StateResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_state(state_id: Identifier) -> StateResponse:
    pass


@router.put(
    "/states/{state_id}",
    response_model=StateResponse,
    status_code=status.HTTP_200_OK,
)
async def update_state(
    state_id: Identifier,
    updates: UpdateStateRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> StateResponse:
    pass


@router.delete(
    "/states/{state_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_state(
    state_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> Response:
    pass


# Places
@router.post(
    "/places",
    response_model=PlaceResponse,
    status_code=status.HTTP_201_CREATED,
)
async def create_place(
    new_place: NewPlaceRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> PlaceResponse:
    pass


@router.get(
    "/places",
    response_model=List[PlaceResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_places() -> List[PlaceResponse]:
    pass


@router.get(
    "/places/{place_id}",
    response_model=PlaceResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_place(place_id: Identifier) -> PlaceResponse:
    pass


@router.put(
    "/places/{place_id}",
    response_model=PlaceResponse,
    status_code=status.HTTP_200_OK,
)
async def update_place(
    place_id: Identifier,
    updates: UpdatePlaceRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> PlaceResponse:
    pass


@router.delete(
    "/places/{place_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_place(
    place_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
) -> Response:
    pass
