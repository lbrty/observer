from fastapi import APIRouter, Depends, Response
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier
from observer.components.pagination import pagination
from observer.schemas.pagination import Pagination
from observer.schemas.users import (
    CreateUserRequest,
    UpdateUserRequest,
    UserResponse,
    UsersResponse,
)

router = APIRouter(prefix="/users")


@router.post(
    "",
    status_code=status.HTTP_201_CREATED,
    response_model=UsersResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_409_CONFLICT,
    ),
    tags=["admin", "users"],
)
async def create_user(new_user: CreateUserRequest) -> UserResponse:
    pass


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
