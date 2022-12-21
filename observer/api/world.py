from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.common.types import Identifier, Role
from observer.components.auth import RequiresRoles, current_user
from observer.components.services import audit_service, world_service
from observer.entities.base import SomeUser
from observer.schemas.world import (
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
from observer.services.audit_logs import AuditServiceInterface
from observer.services.world import WorldServiceInterface

router = APIRouter(prefix="/world")


# Countries
@router.post(
    "/countries",
    response_model=CountryResponse,
    status_code=status.HTTP_201_CREATED,
)
async def create_country(
    tasks: BackgroundTasks,
    new_country: NewCountryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> CountryResponse:
    country = await world.create_country(new_country)
    tag = "endpoint=create_country"
    audit_log = await world.create_log(
        f"{tag},action=create:country,country_id={country.id},ref_id={user.ref_id}",
        None,
        country.dict(exclude={"id"}),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.to_response(country)


@router.get(
    "/countries",
    response_model=List[CountryResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_countries(world: WorldServiceInterface = Depends(world_service)) -> List[CountryResponse]:
    countries = await world.get_countries()
    return await world.list_to_response(countries)


@router.get(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
)
async def get_country(
    country_id: Identifier,
    world: WorldServiceInterface = Depends(world_service),
) -> CountryResponse:
    country = await world.get_country(country_id)
    return await world.to_response(country)


@router.put(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
)
async def update_country(
    tasks: BackgroundTasks,
    country_id: Identifier,
    updates: UpdateCountryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> CountryResponse:
    country = await world.get_country(country_id)
    updated_country = await world.update_country(country_id, updates)
    tag = "endpoint=update_country"
    audit_log = await world.create_log(
        f"{tag},action=update:country,country_id={updated_country.id},ref_id={user.ref_id}",
        None,
        dict(
            old_country=country.dict(exclude={"id"}),
            new_country=updated_country.dict(exclude={"id"}),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.to_response(updated_country)


@router.delete(
    "/countries/{country_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_country(
    tasks: BackgroundTasks,
    country_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    deleted_country = await world.delete_country(country_id)
    tag = "endpoint=delete_country"
    audit_log = await world.create_log(
        f"{tag},action=delete:country,country_id={deleted_country.id},ref_id={user.ref_id}",
        None,
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


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
