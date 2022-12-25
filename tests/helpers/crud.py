from observer.common.types import Identifier, PlaceType
from observer.context import Context
from observer.entities.idp import Category, NewCategory
from observer.entities.projects import NewProject, Project
from observer.entities.world import (
    Country,
    NewCountry,
    NewPlace,
    NewState,
    Place,
    State,
)


async def create_project(ctx: Context, name: str, description: str) -> Project:
    project = await ctx.projects_repo.create_project(NewProject(name=name, description=description))
    return project


async def create_category(ctx: Context, name: str) -> Category:
    category = await ctx.category_repo.create_category(NewCategory(name=name))
    return category


async def create_country(ctx: Context, name: str, code: str) -> Country:
    country = await ctx.world_repo.create_country(NewCountry(name=name, code=code))
    return country


async def create_state(
    ctx: Context,
    name: str,
    code: str,
    country_id: Identifier,
) -> State:
    state = await ctx.world_repo.create_state(
        NewState(
            name=name,
            code=code,
            country_id=country_id,
        )
    )
    return state


async def create_city(
    ctx: Context,
    name: str,
    code: str,
    country_id: Identifier,
    state_id: Identifier,
) -> Place:
    city = await ctx.world_repo.create_place(
        NewPlace(
            name=name,
            code=code,
            place_type=PlaceType.city,
            country_id=country_id,
            state_id=state_id,
        )
    )
    return city


async def create_town(
    ctx: Context,
    name: str,
    code: str,
    country_id: Identifier,
    state_id: Identifier,
) -> Place:
    town = await ctx.world_repo.create_place(
        NewPlace(
            name=name,
            code=code,
            place_type=PlaceType.town,
            country_id=country_id,
            state_id=state_id,
        )
    )
    return town


async def create_village(
    ctx: Context,
    name: str,
    code: str,
    country_id: Identifier,
    state_id: Identifier,
) -> Place:
    village = await ctx.world_repo.create_place(
        NewPlace(
            name=name,
            code=code,
            place_type=PlaceType.village,
            country_id=country_id,
            state_id=state_id,
        )
    )
    return village
