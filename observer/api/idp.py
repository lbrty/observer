from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from starlette import status

from observer.common.types import Identifier, Role, SomeStr
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    category_service,
    idp_service,
    projects_service,
    world_service,
)
from observer.entities.base import SomeUser
from observer.schemas.idp import (
    CategoryResponse,
    IDPResponse,
    NewCategoryRequest,
    NewIDPRequest,
    UpdateCategoryRequest,
)
from observer.services.audit_logs import AuditServiceInterface
from observer.services.categories import CategoryServiceInterface
from observer.services.idp import IDPServiceInterface
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_category,action=create:category",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> CategoryResponse:
    category = await categories.create_category(new_category)
    audit_log = props.new_event(
        f"category_id={category.id},ref_id={user.ref_id}",
        category,
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
    name: SomeStr = Query(None, description="Lookup by name"),
    categories: CategoryServiceInterface = Depends(category_service),
) -> List[CategoryResponse]:
    category_list = await categories.get_categories(name)
    return [CategoryResponse(**category.dict()) for category in category_list]


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
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_category,action=update:category",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> CategoryResponse:
    category = await categories.get_category(category_id)
    updated_category = await categories.update_category(category_id, updates)
    audit_log = props.new_event(
        f"category_id={category.id},ref_id={user.ref_id}",
        dict(
            old_category=category,
            new_category=updated_category,
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
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_category,action=delete:category",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    category = await categories.delete_category(category_id)
    audit_log = props.new_event(
        f"category_id={category.id},ref_id={user.ref_id}",
        category,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/people",
    response_model=IDPResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["idp", "people"],
)
async def create_idp(
    tasks: BackgroundTasks,
    new_idp: NewIDPRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant]),
    ),
    categories: CategoryServiceInterface = Depends(category_service),
    world: WorldServiceInterface = Depends(world_service),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
    idp: IDPServiceInterface = Depends(idp_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_idp,action=create:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> IDPResponse:
    # TODO: add verifications for foreign keys
    person = await idp.create_idp(new_idp)
    audit_log = props.new_event(f"person_id={person.id},ref_id={user.ref_id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return IDPResponse(**person.dict())
