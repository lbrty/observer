from fastapi import APIRouter, Response
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.schemas.users import UserInviteRequest, UserInviteResponse

router = APIRouter(prefix="/invites")


@router.post(
    "",
    response_model=UserInviteResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["admin", "invites"],
)
async def create_invite(
    invite_request: UserInviteRequest,
) -> UserInviteResponse:
    ...


@router.delete(
    "{invite_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["admin", "invites"],
)
async def delete_invite() -> Response:
    ...
