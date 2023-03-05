from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.pagination import pagination
from observer.components.services import audit_service, auth_service, users_service
from observer.entities.users import User, UserUpdate
from observer.schemas.auth import RegistrationPayload
from observer.schemas.pagination import Pagination
from observer.schemas.users import (
    CreateUserRequest,
    UpdateUserRequest,
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
    tags=["admin", "users"],
)
async def get_users(pages: Pagination = Depends(pagination)) -> UsersResponse:
    pass


@router.get(
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
async def get_user(user_id: Identifier) -> UserResponse:
    pass


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
async def update_user(user_id: Identifier, updates: UpdateUserRequest) -> UserResponse:
    pass


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
async def delete_user(user_id: Identifier) -> Response:
    pass
