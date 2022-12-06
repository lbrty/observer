from observer.context import ctx
from observer.services.mfa import MFAServiceInterface


async def mfa_service() -> MFAServiceInterface:
    return ctx.mfa_service
