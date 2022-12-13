from fastapi import Depends

from observer.api.exceptions import TOTPExistsError
from observer.components.auth import authenticated_user
from observer.context import ctx
from observer.entities.users import User
from observer.services.mfa import MFAServiceInterface


async def mfa_service() -> MFAServiceInterface:
    if ctx.mfa_service:
        return ctx.mfa_service

    raise RuntimeError("MFAService is None")


async def user_with_no_mfa(user: User = Depends(authenticated_user)):
    """Checks if user has empty MFA configuration.

    If MFA configuration exists then it raises `TOTPException`.

    Args:
        user(User): injected user instance

    Returns:
        user: User instance
    """
    if user.mfa_enabled and user.mfa_encrypted_secret is not None:
        raise TOTPExistsError

    return user
