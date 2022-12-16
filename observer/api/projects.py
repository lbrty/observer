from fastapi import APIRouter, BackgroundTasks, Depends
from starlette import status

from observer.common.permissions import permission_matrix
from observer.common.types import Role
from observer.components.auth import RequiresRoles
from observer.components.projects import project_details
from observer.components.services import (
    audit_service,
    permissions_service,
    projects_service,
)
from observer.entities.base import SomeUser
from observer.entities.projects import Project
from observer.schemas.permissions import NewPermission
from observer.schemas.projects import NewProjectRequest, ProjectResponse
from observer.services.audit_logs import AuditServiceInterface
from observer.services.permissions import PermissionsServiceInterface
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
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> ProjectResponse:
    project = await projects.create_project(new_project)
    tag = "endpoint=create_project"
    audit_log = await projects.create_log(
        f"{tag},action=create:project,project_id={project.id},ref_id={user.ref_id}",
        None,
        dict(
            id=str(project.id),
            name=project.name,
            description=project.description,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)

    base_permission = permission_matrix[user.role]
    permission = await permissions.create_permission(
        NewPermission(
            **base_permission.dict(),
            user_id=user.id,
            project_id=project.id,
        )
    )
    audit_log = await projects.create_log(
        f"{tag},action=create:permission,permission_id={permission.id},ref_id={user.ref_id}",
        None,
        dict(
            project_id=str(project.id),
            project_name=str(project.name),
            user_id=str(user.id),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(project)


@router.get("/{project_id}", response_model=ProjectResponse, status_code=status.HTTP_200_OK)
async def get_project(project: Project = Depends(project_details)) -> ProjectResponse:
    return ProjectResponse(**project.dict())
