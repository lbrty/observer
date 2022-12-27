from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import Identifier, PlaceFilters, Role, StateFilters
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles, current_user
from observer.components.filters import place_filters, state_filters
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_country,action=create:country",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> CountryResponse:
    country = await world.create_country(new_country)
    audit_log = props.new_event(
        f"country_id={country.id},ref_id={user.ref_id}",
        jsonable_encoder(country.dict(exclude={"id"})),
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_country,action=update:country",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> CountryResponse:
    country = await world.get_country(country_id)
    updated_country = await world.update_country(country_id, updates)
    audit_log = props.new_event(
        f"country_id={country.id},ref_id={user.ref_id}",
        jsonable_encoder(
            dict(
                old_country=country.dict(exclude={"id"}),
                new_country=updated_country.dict(exclude={"id"}),
            )
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_country,action=delete:country",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    deleted_country = await world.delete_country(country_id)
    audit_log = props.new_event(
        f"country_id={deleted_country.id},ref_id={user.ref_id}",
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_state,action=create:state",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> StateResponse:
    state = await world.create_state(new_state)
    audit_log = props.new_event(
        f"state_id={state.id},ref_id={user.ref_id}",
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
async def get_states(
    filters: StateFilters = Depends(state_filters),
    world: WorldServiceInterface = Depends(world_service),
) -> List[StateResponse]:
    states = await world.get_states(filters)
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_state,action=update:state",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> StateResponse:
    state = await world.get_state(state_id)
    updated_state = await world.update_state(state_id, updates)
    audit_log = props.new_event(
        f"state_id={updated_state.id},ref_id={user.ref_id}",
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_state,action=delete:state",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    deleted_state = await world.delete_state(state_id)
    audit_log = props.new_event(
        f"state_id={deleted_state.id},ref_id={user.ref_id}",
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_place,action=create:place",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> PlaceResponse:
    place = await world.create_place(new_place)
    audit_log = props.new_event(
        f"place_id={place.id},ref_id={user.ref_id}",
        jsonable_encoder(place),
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
async def get_places(
    filters: PlaceFilters = Depends(place_filters),
    world: WorldServiceInterface = Depends(world_service),
) -> List[PlaceResponse]:
    places = await world.get_places(filters)
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_place,action=update:place",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> PlaceResponse:
    await world.get_country(updates.country_id)
    await world.get_state(updates.state_id)

    place = await world.get_place(place_id)
    updated_place = await world.update_place(place_id, updates)
    audit_log = props.new_event(
        f"place_id={updated_place.id},ref_id={user.ref_id}",
        dict(
            old_place=jsonable_encoder(place),
            new_place=jsonable_encoder(updated_place),
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_place,action=delete:place",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    deleted_place = await world.delete_place(place_id)
    audit_log = props.new_event(
        f"place_id={deleted_place.id},ref_id={user.ref_id}",
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
