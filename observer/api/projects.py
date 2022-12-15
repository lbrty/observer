from fastapi import APIRouter, BackgroundTasks, Depends
from starlette import status

from observer.common.types import Role
from observer.components.auth import RequiresRoles
from observer.components.projects import project_details
from observer.components.services import audit_service, projects_service
from observer.entities.base import SomeUser
from observer.entities.projects import Project
from observer.schemas.projects import NewProjectRequest, ProjectResponse
from observer.services.audit_logs import AuditServiceInterface
from observer.services.projects import ProjectsServiceInterface

router = APIRouter(prefix="/projects")


@router.post("/", response_model=ProjectResponse, status_code=status.HTTP_201_CREATED)
async def create_project(
    tasks: BackgroundTasks,
    new_project: NewProjectRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant]),
    ),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> ProjectResponse:
    project = await projects.create_project(new_project)
    audit_log = await projects.create_log(
        f"endpoint=create_project,action=create:project,project_id={project.id},ref_id={user.ref_id}",
        None,
        dict(
            id=str(project.id),
            name=project.name,
            description=project.description,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(project)


@router.get("/{project_id}", response_model=ProjectResponse, status_code=status.HTTP_200_OK)
async def get_project(project: Project = Depends(project_details)) -> ProjectResponse:
    return ProjectResponse(**project.dict())
