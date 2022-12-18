from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.api.exceptions import ConflictError, NotFoundError
from observer.common.exceptions import get_api_errors
from observer.common.permissions import permission_matrix
from observer.common.types import Role
from observer.components.auth import RequiresRoles, current_user
from observer.components.pagination import pagination
from observer.components.projects import (
    deletable_project,
    invitable_project,
    owned_project,
    updatable_project,
    viewable_project,
)
from observer.components.services import (
    audit_service,
    permissions_service,
    projects_service,
    users_service,
)
from observer.entities.base import SomeUser
from observer.entities.projects import Project
from observer.entities.users import User
from observer.schemas.pagination import Pagination
from observer.schemas.permissions import (
    BasePermission,
    NewPermissionRequest,
    UpdatePermissionRequest,
)
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
from observer.services.users import UsersServiceInterface

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
    new_project.owner_id = str(user.id)
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

    # According to role we need to use default permissions
    base_permission = permission_matrix[user.role]
    permission = await permissions.create_permission(
        NewPermissionRequest(
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
    project: Project = Depends(viewable_project),
) -> ProjectResponse:
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
    user: User = Depends(current_user),
    project: Project = Depends(updatable_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> ProjectResponse:
    tag = "endpoint=update_project"
    updated_project = await projects.update_project(project.id, updates)
    audit_log = await projects.create_log(
        f"{tag},action=update:project,project_id={project.id},ref_id={user.ref_id}",
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
    user: User = Depends(current_user),
    project: Project = Depends(owned_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_project"

    await projects.delete_project(project.id)
    audit_log = await projects.create_log(
        f"{tag},action=delete:project,project_id={project.id},ref_id={user.ref_id}",
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
    project: Project = Depends(viewable_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
    pages: Pagination = Depends(pagination),
) -> ProjectMembersResponse:
    members = await projects.get_members(project.id, pages.offset, pages.limit)
    return ProjectMembersResponse(items=[ProjectMemberResponse(**member.dict()) for member in members])


@router.post(
    "/{project_id}/members",
    response_model=ProjectMemberResponse,
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_200_OK,
)
async def add_project_member(
    tasks: BackgroundTasks,
    new_permission: NewPermissionRequest,
    user: User = Depends(current_user),
    project: Project = Depends(invitable_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
    users: UsersServiceInterface = Depends(users_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> ProjectMemberResponse:
    tag = "endpoint=add_project_member"
    member_user = await users.get_by_id(new_permission.user_id)
    if not member_user:
        raise NotFoundError(message="User not found")

    if str(project.id) != str(new_permission.project_id):
        raise ConflictError(message="Project ids in path and payload differ")

    permission = await permissions.create_permission(
        NewPermissionRequest(**new_permission.dict()),
    )
    audit_log = await projects.create_log(
        f"{tag},action=create:permission,permission_id={permission.id},ref_id={user.ref_id}",
        None,
        dict(member_ref_id=str(member_user.ref_id)),
    )
    tasks.add_task(audits.add_event, audit_log)
    return ProjectMemberResponse(**dict(**member_user.dict(), permissions=BasePermission(**permission.dict())))


@router.put(
    "/{project_id}/members",
    response_model=ProjectMemberResponse,
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_200_OK,
)
async def update_project_members_permissions(
    permission: UpdatePermissionRequest,
    project: Project = Depends(invitable_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
) -> ProjectMemberResponse:
    pass


@router.delete(
    "/{project_id}/members/{user_id}",
    responses=get_api_errors(status.HTTP_404_NOT_FOUND, status.HTTP_403_FORBIDDEN),
    status_code=status.HTTP_204_NO_CONTENT,
)
async def delete_project_member(
    tasks: BackgroundTasks,
    user: User = Depends(current_user),
    project: Project = Depends(deletable_project),
    projects: ProjectsServiceInterface = Depends(projects_service),
    audits: AuditServiceInterface = Depends(audit_service),
) -> Response:
    tag = "endpoint=delete_project_member"
    deleted_permission = await projects.delete_member(project.id, user.id)
    audit_log = await projects.create_log(
        f"{tag},action=delete:permission,permission_id={deleted_permission.id},ref_id={user.ref_id}",
        None,
        deleted_permission.dict(),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
