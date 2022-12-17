from typing import Tuple

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.permissions import permission_matrix
from observer.common.types import Role
from observer.components.auth import RequiresRoles
from observer.components.pagination import pagination
from observer.components.projects import project_details, project_with_member
from observer.components.services import (
    audit_service,
    permissions_service,
    projects_service,
)
from observer.entities.base import SomeUser
from observer.entities.projects import Project, ProjectMember
from observer.schemas.pagination import Pagination
from observer.schemas.permissions import NewPermission
from observer.schemas.projects import (
    NewProjectRequest,
    ProjectMemberResponse,
    ProjectMembersResponse,
    ProjectResponse,
    UpdateProjectRequest,
)
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
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(project)


@router.get(
    "/{project_id}",
    response_model=ProjectResponse,
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_200_OK,
)
async def get_project(
    project_and_member: Tuple[Project, ProjectMember] = Depends(project_with_member),
) -> ProjectResponse:
    project, _ = project_and_member
    return ProjectResponse(**project.dict())


@router.put(
    "/{project_id}",
    response_model=ProjectResponse,
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_200_OK,
)
async def update_project(
    tasks: BackgroundTasks,
    updates: UpdateProjectRequest,
    project_and_member: Tuple[Project, ProjectMember] = Depends(project_with_member),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> ProjectResponse:
    tag = "endpoint=update_project"
    project, member = project_and_member
    updated_project = await projects.update_project(project.id, updates)
    audit_log = await projects.create_log(
        f"{tag},action=update:project,project_id={project.id},ref_id={member.ref_id}",
        None,
        dict(
            old_project=project.dict(exclude={"id"}),
            new_project=updated_project.dict(exclude={"id"}),
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(updated_project)


@router.delete(
    "/{project_id}",
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_project(
    tasks: BackgroundTasks,
    project_and_member: Tuple[Project, ProjectMember] = Depends(project_with_member),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_project"
    project, member = project_and_member
    await projects.delete_project(project.id)
    audit_log = await projects.create_log(
        f"{tag},action=delete:project,project_id={project.id},ref_id={member.ref_id}",
        None,
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.get(
    "/{project_id}/members",
    response_model=ProjectMembersResponse,
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_200_OK,
)
async def get_project_members(
    project: Project = Depends(project_details),
    projects: ProjectsServiceInterface = Depends(projects_service),
    pages: Pagination = Depends(pagination),
) -> ProjectMembersResponse:
    members = await projects.get_members(project.id, pages.offset, pages.limit)
    return ProjectMembersResponse(items=[ProjectMemberResponse(**member.dict()) for member in members])
