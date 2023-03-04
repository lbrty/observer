from fastapi import APIRouter

from observer.api.admin import invites, users

router = APIRouter(prefix="/admin")

router.include_router(invites.router)
router.include_router(users.router)
