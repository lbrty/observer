from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import Identifier, Role, SomeStr
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    category_service,
    projects_service,
    world_service,
)
from observer.entities.base import SomeUser
from observer.schemas.displaced_persons import (
    CategoryResponse,
    DisplacedPersonResponse,
    NewCategoryRequest,
    UpdateCategoryRequest,
)
from observer.services.audit_logs import AuditServiceInterface
from observer.services.categories import CategoryServiceInterface
from observer.services.projects import ProjectsServiceInterface
from observer.services.world import WorldServiceInterface

router = APIRouter(prefix="/idp")


@router.post(
    "/categories",
    response_model=CategoryResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["idp", "categories"],
)
async def create_category(
    tasks: BackgroundTasks,
    new_category: NewCategoryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: CategoryServiceInterface = Depends(category_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> CategoryResponse:
    tag = "endpoint=create_category"
    category = await categories.create_category(new_category)
    audit_log = await categories.create_log(
        f"{tag},action=create:category,category_id={category.id},ref_id={user.ref_id}",
        None,
        jsonable_encoder(category),
    )
    tasks.add_task(audits.add_event, audit_log)
    return CategoryResponse(**category.dict())


@router.get(
    "/categories",
    response_model=List[CategoryResponse],
    status_code=status.HTTP_200_OK,
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["idp", "categories"],
)
async def get_categories(
    name: SomeStr = Query(..., description="Lookup by name"),
    categories: CategoryServiceInterface = Depends(category_service),
) -> List[CategoryResponse]:
    categories = await categories.get_categories(name)
    return [CategoryResponse(**category.dict()) for category in categories]


@router.get(
    "/categories/{category_id}",
    response_model=CategoryResponse,
    status_code=status.HTTP_200_OK,
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["idp", "categories"],
)
async def get_category(
    category_id: Identifier,
    categories: CategoryServiceInterface = Depends(category_service),
) -> CategoryResponse:
    category = await categories.get_category(category_id)
    return CategoryResponse(**category.dict())


@router.put(
    "/categories/{category_id}",
    response_model=CategoryResponse,
    status_code=status.HTTP_200_OK,
    tags=["idp", "categories"],
)
async def update_category(
    tasks: BackgroundTasks,
    category_id: Identifier,
    updates: UpdateCategoryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: CategoryServiceInterface = Depends(category_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> CategoryResponse:
    tag = "endpoint=update_category"
    category = await categories.get_category(category_id)
    updated_category = await categories.update_category(category_id, updates)
    audit_log = await categories.create_log(
        f"{tag},action=update:category,category_id={category_id},ref_id={user.ref_id}",
        None,
        dict(
            old_category=jsonable_encoder(category),
            new_category=jsonable_encoder(updated_category),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return CategoryResponse(**category.dict())


@router.delete(
    "/categories/{category_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["idp", "categories"],
)
async def delete_category(
    tasks: BackgroundTasks,
    category_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: CategoryServiceInterface = Depends(category_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_category"
    category = await categories.delete_category(category_id)
    audit_log = await categories.create_log(
        f"{tag},action=delete:category,category_id={category_id},ref_id={user.ref_id}",
        None,
        jsonable_encoder(category),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/people",
    response_model=DisplacedPersonResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["idp", "people"],
)
async def create_idp(
    tasks: BackgroundTasks,
    new_category: NewCategoryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: CategoryServiceInterface = Depends(category_service),
    world: WorldServiceInterface = Depends(world_service),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> CategoryResponse:
    tag = "endpoint=create_idp"
    category = await categories.create_category(new_category)
    audit_log = await categories.create_log(
        f"{tag},action=create:category,category_id={category.id},ref_id={user.ref_id}",
        None,
        jsonable_encoder(category),
    )
    tasks.add_task(audits.add_event, audit_log)
    return CategoryResponse(**category.dict())
