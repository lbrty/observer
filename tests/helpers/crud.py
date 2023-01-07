from observer.common.types import Identifier, PetStatus, PlaceType
from observer.context import Context
from observer.entities.idp import IDP, Category, NewCategory, NewIDP
from observer.entities.permissions import NewPermission, Permission
from observer.entities.pets import NewPet, Pet
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
    project = await ctx.repos.projects.create_project(NewProject(name=name, description=description))
    return project


async def create_permission(ctx: Context, new_permission: NewPermission) -> Permission:
    permission = await ctx.repos.permissions.create_permission(new_permission)
    return permission


async def create_category(ctx: Context, name: str) -> Category:
    category = await ctx.repos.category.create_category(NewCategory(name=name))
    return category


async def create_country(ctx: Context, name: str, code: str) -> Country:
    country = await ctx.repos.world.create_country(NewCountry(name=name, code=code))
    return country


async def create_state(
    ctx: Context,
    name: str,
    code: str,
    country_id: Identifier,
) -> State:
    state = await ctx.repos.world.create_state(
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
    city = await ctx.repos.world.create_place(
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
    town = await ctx.repos.world.create_place(
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
    village = await ctx.repos.world.create_place(
        NewPlace(
            name=name,
            code=code,
            place_type=PlaceType.village,
            country_id=country_id,
            state_id=state_id,
        )
    )
    return village


async def create_person(
    ctx: Context,
    project_id: Identifier,
) -> IDP:
    person = await ctx.repos.idp.create_idp(
        NewIDP(
            project_id=project_id,
            email="Full_Name@examples.com",
            full_name="Full Name",
            phone_number="+11111111",
            tags=["one", "two"],
        )
    )
    return person


async def create_pet(
    ctx: Context,
    name: str,
    status: PetStatus,
    registration_id: Identifier,
    project_id: Identifier,
    owner_id: Identifier,
) -> Pet:
    pet = await ctx.repos.pets.create_pet(
        NewPet(
            name=name,
            notes="Petya",
            status=status,
            registration_id=registration_id,
            owner_id=owner_id,
            project_id=project_id,
        )
    )
    return pet
