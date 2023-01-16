from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Role, SomeStr
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import audit_service, category_service
from observer.entities.base import SomeUser
from observer.schemas.people import (
    CategoryResponse,
    NewCategoryRequest,
    UpdateCategoryRequest,
)
from observer.services.audit_logs import IAuditService
from observer.services.categories import ICategoryService

router = APIRouter(prefix="/categories")


@router.post(
    "",
    response_model=CategoryResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    tags=["people", "categories"],
)
async def create_category(
    tasks: BackgroundTasks,
    new_category: NewCategoryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: ICategoryService = Depends(category_service),
    audits: IAuditService = Depends(audit_service),
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
    "",
    response_model=List[CategoryResponse],
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["people", "categories"],
)
async def get_categories(
    name: SomeStr = Query(None, description="Lookup by name"),
    categories: ICategoryService = Depends(category_service),
) -> List[CategoryResponse]:
    category_list = await categories.get_categories(name)
    return [CategoryResponse(**category.dict()) for category in category_list]


@router.get(
    "/{category_id}",
    response_model=CategoryResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["people", "categories"],
)
async def get_category(
    category_id: Identifier,
    categories: ICategoryService = Depends(category_service),
) -> CategoryResponse:
    category = await categories.get_category(category_id)
    return CategoryResponse(**category.dict())


@router.put(
    "/{category_id}",
    response_model=CategoryResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["people", "categories"],
)
async def update_category(
    tasks: BackgroundTasks,
    category_id: Identifier,
    updates: UpdateCategoryRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: ICategoryService = Depends(category_service),
    audits: IAuditService = Depends(audit_service),
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
    "/{category_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["people", "categories"],
)
async def delete_category(
    tasks: BackgroundTasks,
    category_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    categories: ICategoryService = Depends(category_service),
    audits: IAuditService = Depends(audit_service),
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
