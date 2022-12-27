from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.permissions import (
    assert_can_see_private_info,
    assert_deletable,
    assert_updatable,
    assert_viewable,
    assert_writable,
)
from observer.common.types import Identifier, Role, SomeStr
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles, authenticated_user
from observer.components.services import (
    audit_service,
    category_service,
    idp_service,
    permissions_service,
    secrets_service,
)
from observer.entities.base import SomeUser
from observer.entities.idp import PersonalInfo
from observer.schemas.idp import (
    CategoryResponse,
    IDPResponse,
    NewCategoryRequest,
    NewIDPRequest,
    PersonalInfoResponse,
    UpdateCategoryRequest,
    UpdateIDPRequest,
)
from observer.services.audit_logs import AuditServiceInterface
from observer.services.categories import CategoryServiceInterface
from observer.services.idp import IDPServiceInterface
from observer.services.permissions import PermissionsServiceInterface
from observer.services.secrets import SecretsServiceInterface

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
    user: SomeUser = Depends(authenticated_user),
    audits: AuditServiceInterface = Depends(audit_service),
    idp: IDPServiceInterface = Depends(idp_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_idp,action=create:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> IDPResponse:
    permission = await permissions.find(new_idp.project_id, user.id)
    assert_writable(user, permission)

    person = await idp.create_idp(new_idp)
    audit_log = props.new_event(f"person_id={person.id},ref_id={user.ref_id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return IDPResponse(**person.dict())


@router.get(
    "/people/{idp_id}",
    response_model=IDPResponse,
    status_code=status.HTTP_200_OK,
    tags=["idp", "people"],
)
async def get_idp(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IDPServiceInterface = Depends(idp_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    secrets: SecretsServiceInterface = Depends(secrets_service),
) -> IDPResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    idp_record = await secrets.anonymize_idp(idp_record)
    return IDPResponse(**idp_record.dict())


@router.get(
    "/people/{idp_id}/personal-info",
    response_model=PersonalInfoResponse,
    response_model_exclude_none=True,
    status_code=status.HTTP_200_OK,
    tags=["idp", "people"],
)
async def get_personal_info(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IDPServiceInterface = Depends(idp_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    secrets: SecretsServiceInterface = Depends(secrets_service),
) -> PersonalInfoResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_can_see_private_info(user, permission)
    pi = await secrets.decrypt_personal_info(
        PersonalInfo(
            email=idp_record.email,
            phone_number=idp_record.phone_number,
            phone_number_additional=idp_record.phone_number_additional,
        ),
    )
    pi.full_name = idp_record.full_name
    return PersonalInfoResponse(**pi.dict())


@router.put(
    "/people/{idp_id}",
    response_model=IDPResponse,
    status_code=status.HTTP_200_OK,
    tags=["idp", "people"],
)
async def update_idp(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    idp_updates: UpdateIDPRequest,
    user: SomeUser = Depends(authenticated_user),
    idp: IDPServiceInterface = Depends(idp_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    audits: AuditServiceInterface = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_idp,action=update:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> IDPResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_updatable(user, permission)
    updated_idp = await idp.update_idp(idp_id, idp_updates)
    audit_log = props.new_event(
        f"person_id={updated_idp.id},ref_id={user.ref_id}",
        jsonable_encoder(
            updated_idp.dict(
                exclude_none=True,
                exclude={"email", "phone_number", "phone_number_additional"},
            )
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return IDPResponse(**updated_idp.dict())


@router.delete(
    "/people/{idp_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["idp", "people"],
)
async def delete_idp(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IDPServiceInterface = Depends(idp_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    audits: AuditServiceInterface = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_idp,action=delete:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_deletable(user, permission)
    deleted_idp = await idp.delete_idp(idp_id)
    audit_log = props.new_event(
        f"person_id={deleted_idp.id},ref_id={user.ref_id}",
        jsonable_encoder(deleted_idp.dict(exclude_none=True)),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
