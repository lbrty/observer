from fastapi import APIRouter
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.schemas.auth import TokenResponse
from observer.schemas.users import InviteJoinRequest

router = APIRouter(prefix="/invites")


@router.post(
    "/join/{code}",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["invites"],
)
async def join_with_invite(join_request: InviteJoinRequest) -> TokenResponse:
    ...
