from fastapi import APIRouter, Depends
from starlette import status

from observer.components.auth import authenticated_user
from observer.entities.users import User

router = APIRouter(prefix="/account")


@router.get("/confirmation/{code}", status_code=status.HTTP_200_OK)
async def confirm_account(code: str):
    pass


@router.get("/confirmation/resend", status_code=status.HTTP_204_NO_CONTENT)
async def resend_confirmation(user: User = Depends(authenticated_user)):
    pass
