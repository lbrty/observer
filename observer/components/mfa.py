from fastapi import Depends

from api.exceptions import TOTPExistsError
from components.auth import authenticated_user
from entities.users import User
from observer.context import ctx
from observer.services.mfa import MFAServiceInterface


async def mfa_service() -> MFAServiceInterface:
    return ctx.mfa_service


async def user_with_no_mfa(user: User = Depends(authenticated_user)):
    if user.mfa_enabled and user.mfa_encrypted_secret is not None:
        raise TOTPExistsError
