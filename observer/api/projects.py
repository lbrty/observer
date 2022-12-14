from fastapi import APIRouter, Depends
from starlette import status

from observer.common.types import Identifier
from observer.components.auth import current_user
from observer.entities.users import User
from observer.schemas.projects import ProjectResponse

router = APIRouter(prefix="/projects")


@router.get("/{project_id}", response_model=ProjectResponse, status_code=status.HTTP_200_OK)
async def get_project(project_id: Identifier, user: User = Depends(current_user)) -> ProjectResponse:
    pass
