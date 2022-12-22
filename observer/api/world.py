from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.api.exceptions import NotFoundError
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
    tags=["world", "places"],
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
    tag = "endpoint=create_country"
    country = await world.create_country(new_country)
    audit_log = await world.create_log(
        f"{tag},action=create:country,country_id={country.id},ref_id={user.ref_id}",
        None,
        country.dict(exclude={"id"}),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.country_to_response(country)


@router.get(
    "/countries",
    response_model=List[CountryResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "places"],
)
async def get_countries(world: WorldServiceInterface = Depends(world_service)) -> List[CountryResponse]:
    countries = await world.get_countries()
    return await world.countries_to_response(countries)


@router.get(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "places"],
)
async def get_country(
    country_id: Identifier,
    world: WorldServiceInterface = Depends(world_service),
) -> CountryResponse:
    country = await world.get_country(country_id)
    return await world.country_to_response(country)


@router.put(
    "/countries/{country_id}",
    response_model=CountryResponse,
    status_code=status.HTTP_200_OK,
    tags=["world", "places"],
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
    if not country:
        raise NotFoundError(message="Country not found")

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
    return await world.country_to_response(updated_country)


@router.delete(
    "/countries/{country_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["world", "places"],
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
    tag = "endpoint=delete_country"
    deleted_country = await world.delete_country(country_id)
    if not deleted_country:
        raise NotFoundError(message="Country not found")

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
    tags=["world", "states"],
)
async def create_state(
    tasks: BackgroundTasks,
    new_state: NewStateRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> StateResponse:
    tag = "endpoint=create_state"
    state = await world.create_state(new_state)
    audit_log = await world.create_log(
        f"{tag},action=create:state,state_id={state.id},ref_id={user.ref_id}",
        None,
        jsonable_encoder(state),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.state_to_response(state)


@router.get(
    "/states",
    response_model=List[StateResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "states"],
)
async def get_states(world: WorldServiceInterface = Depends(world_service)) -> List[StateResponse]:
    states = await world.get_states()
    return await world.states_to_response(states)


@router.get(
    "/states/{state_id}",
    response_model=StateResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "states"],
)
async def get_state(
    state_id: Identifier,
    world: WorldServiceInterface = Depends(world_service),
) -> StateResponse:
    state = await world.get_state(state_id)
    return await world.state_to_response(state)


@router.put(
    "/states/{state_id}",
    response_model=StateResponse,
    status_code=status.HTTP_200_OK,
    tags=["world", "states"],
)
async def update_state(
    tasks: BackgroundTasks,
    state_id: Identifier,
    updates: UpdateStateRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> StateResponse:
    state = await world.get_state(state_id)
    if not state:
        raise NotFoundError(message="State not found")

    updated_state = await world.update_state(state_id, updates)
    tag = "endpoint=update_state"
    audit_log = await world.create_log(
        f"{tag},action=update:state,state_id={updated_state.id},ref_id={user.ref_id}",
        None,
        dict(
            old_state=jsonable_encoder(state),
            new_state=jsonable_encoder(updated_state),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.state_to_response(updated_state)


@router.delete(
    "/states/{state_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["world", "states"],
)
async def delete_state(
    tasks: BackgroundTasks,
    state_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_state"
    deleted_state = await world.delete_state(state_id)
    if not deleted_state:
        raise NotFoundError(message="State not found")

    audit_log = await world.create_log(
        f"{tag},action=delete:state,state_id={deleted_state.id},ref_id={user.ref_id}",
        None,
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


# Places
@router.post(
    "/places",
    response_model=PlaceResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["world", "places"],
)
async def create_place(
    tasks: BackgroundTasks,
    new_place: NewPlaceRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> PlaceResponse:
    tag = "endpoint=create_place"
    place = await world.create_place(new_place)
    audit_log = await world.create_log(
        f"{tag},action=create:place,place_id={place.id},ref_id={user.ref_id}",
        None,
        place.dict(exclude={"id"}),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.place_to_response(place)


@router.get(
    "/places",
    response_model=List[PlaceResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "places"],
)
async def get_places(world: WorldServiceInterface = Depends(world_service)) -> List[PlaceResponse]:
    places = await world.get_places()
    return await world.places_to_response(places)


@router.get(
    "/places/{place_id}",
    response_model=PlaceResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_user)],
    tags=["world", "places"],
)
async def get_place(
    place_id: Identifier,
    world: WorldServiceInterface = Depends(world_service),
) -> PlaceResponse:
    place = await world.get_place(place_id)
    return await world.place_to_response(place)


@router.put(
    "/places/{place_id}",
    response_model=PlaceResponse,
    status_code=status.HTTP_200_OK,
    tags=["world", "places"],
)
async def update_place(
    tasks: BackgroundTasks,
    place_id: Identifier,
    updates: UpdatePlaceRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> PlaceResponse:
    place = await world.get_place(place_id)
    if not place:
        raise NotFoundError(message="Place not found")

    updated_place = await world.update_place(place_id, updates)
    tag = "endpoint=update_place"
    audit_log = await world.create_log(
        f"{tag},action=update:place,place_id={updated_place.id},ref_id={user.ref_id}",
        None,
        dict(
            old_place=place.dict(exclude={"id"}),
            new_place=updated_place.dict(exclude={"id"}),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await world.place_to_response(updated_place)


@router.delete(
    "/places/{place_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["world", "places"],
)
async def delete_place(
    tasks: BackgroundTasks,
    place_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    world: WorldServiceInterface = Depends(world_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_place"
    deleted_place = await world.delete_place(place_id)
    if not deleted_place:
        raise NotFoundError(message="Place not found")

    audit_log = await world.create_log(
        f"{tag},action=delete:place,place_id={deleted_place.id},ref_id={user.ref_id}",
        None,
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
