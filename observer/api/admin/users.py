from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Pagination, Role, UserFilters
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.filters import user_filters
from observer.components.pagination import pagination
from observer.components.services import audit_service, auth_service, users_service
from observer.entities.users import User, UserUpdate
from observer.schemas.auth import RegistrationPayload
from observer.schemas.users import (
    AdminUpdateUserRequest,
    CreateUserRequest,
    UserResponse,
    UsersResponse,
)
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.users import IUsersService

router = APIRouter(prefix="/users")


@router.post(
    "",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_409_CONFLICT,
    ),
    tags=["admin", "users"],
)
async def admin_create_user(
    new_user: CreateUserRequest,
    tasks: BackgroundTasks,
    audits: IAuditService = Depends(audit_service),
    users: IUsersService = Depends(users_service),
    auth: IAuthService = Depends(auth_service),
    user: User = Depends(RequiresRoles([Role.admin])),
    props: Props = Depends(
        Tracked(
            tag="endpoint=admin_create_user,action=create:user",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    created_user, _ = await auth.register(
        RegistrationPayload(
            email=new_user.email,
            password=new_user.password.get_secret_value(),
        )
    )
    await users.update_user(created_user.id, UserUpdate(**new_user.dict()))
    await users.just_confirm_user(created_user.id)
    audit_log = props.new_event(f"ref_id={user.id}", data=dict(new_user_id=created_user.id))
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.get(
    "",
    status_code=status.HTTP_200_OK,
    response_model=UsersResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    dependencies=[Depends(RequiresRoles([Role.admin]))],
    tags=["admin", "users"],
)
async def admin_get_users(
    users: IUsersService = Depends(users_service),
    filters: UserFilters = Depends(user_filters),
    pages: Pagination = Depends(pagination),
) -> UsersResponse:
    count, items = await users.filter_users(filters, pages)
    return UsersResponse(total=count, items=[UserResponse(**user.dict()) for user in items])


@router.get(
    "/{user_id}",
    status_code=status.HTTP_200_OK,
    response_model=UserResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[Depends(RequiresRoles([Role.admin]))],
    tags=["admin", "users"],
)
async def admin_get_user(
    user_id: Identifier,
    users: IUsersService = Depends(users_service),
) -> UserResponse:
    user = await users.get_by_id(user_id)
    return UserResponse(**user.dict())


@router.put(
    "/{user_id}",
    status_code=status.HTTP_200_OK,
    response_model=UserResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["admin", "users"],
)
async def admin_update_user(
    tasks: BackgroundTasks,
    user_id: Identifier,
    updates: AdminUpdateUserRequest,
    audits: IAuditService = Depends(audit_service),
    users: IUsersService = Depends(users_service),
    user: User = Depends(RequiresRoles([Role.admin])),
    props: Props = Depends(
        Tracked(
            tag="endpoint=admin_update_user,action=update:user",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> UserResponse:
    updated_user = await users.update_user(user_id, UserUpdate(**updates.dict()))
    audit_log = props.new_event(f"ref_id={user.id}", data=jsonable_encoder(updates, exclude_none=True))
    tasks.add_task(audits.add_event, audit_log)
    return UserResponse(**updated_user.dict())


@router.delete(
    "/{user_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["admin", "users"],
)
async def admin_delete_user(
    tasks: BackgroundTasks,
    user_id: Identifier,
    audits: IAuditService = Depends(audit_service),
    users: IUsersService = Depends(users_service),
    user: User = Depends(RequiresRoles([Role.admin])),
    props: Props = Depends(
        Tracked(
            tag="endpoint=admin_delete_user,action=delete:user",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    await users.delete_user(user_id)
    audit_log = props.new_event(f"ref_id={user.id}", data=dict(user_id=user_id))
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
