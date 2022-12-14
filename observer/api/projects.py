from fastapi import APIRouter, Depends
from starlette import status

from observer.components.projects import project_details
from observer.entities.projects import Project
from observer.schemas.projects import ProjectResponse

router = APIRouter(prefix="/projects")


@router.get("/{project_id}", response_model=ProjectResponse, status_code=status.HTTP_200_OK)
async def get_project(project: Project = Depends(project_details)) -> ProjectResponse:
    return ProjectResponse(**project.dict())
