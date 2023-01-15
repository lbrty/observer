from fastapi import APIRouter

from observer.api.admin import invites

router = APIRouter(prefix="/admin")

router.include_router(invites.router)
