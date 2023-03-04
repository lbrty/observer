from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.api.exceptions import ConflictError, NotFoundError
from observer.common.exceptions import get_api_errors
from observer.common.permissions import permission_matrix
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles, authenticated_user, current_user
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
from observer.entities.projects import Project
from observer.entities.users import User
from observer.schemas.pagination import Pagination
from observer.schemas.permissions import (
    NewPermissionRequest,
    PermissionResponse,
    UpdatePermissionRequest,
)
from observer.schemas.projects import (
    NewProjectRequest,
    ProjectMemberResponse,
    ProjectMembersResponse,
    ProjectResponse,
    UpdateProjectRequest,
)
from observer.services.audit_logs import IAuditService
from observer.services.permissions import IPermissionsService
from observer.services.projects import IProjectsService
from observer.services.users import IUsersService

router = APIRouter(prefix="/projects")


@router.post(
    "/",
    response_model=ProjectResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    tags=["projects"],
)
async def create_project(
    tasks: BackgroundTasks,
    new_project: NewProjectRequest,
    user: User = Depends(
        RequiresRoles([Role.admin, Role.consultant]),
    ),
    projects: IProjectsService = Depends(projects_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_project",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> ProjectResponse:
    new_project.owner_id = user.id
    project = await projects.create_project(new_project)
    audit_log = props.new_event(
        f"action=create:project,project_id={project.id},ref_id={user.id}",
        jsonable_encoder(project),
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
    audit_log = props.new_event(
        f"action=create:permission,permission_id={permission.id},ref_id={user.id}",
        jsonable_encoder(project),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(project)


@router.get(
    "/{project_id}",
    response_model=ProjectResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project not found"),
    ),
    status_code=status.HTTP_200_OK,
    tags=["projects"],
)
async def get_project(project: Project = Depends(viewable_project)) -> ProjectResponse:
    return ProjectResponse(**project.dict())


@router.put(
    "/{project_id}",
    response_model=ProjectResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project not found"),
    ),
    status_code=status.HTTP_200_OK,
    tags=["projects"],
)
async def update_project(
    tasks: BackgroundTasks,
    updates: UpdateProjectRequest,
    user: User = Depends(current_user),
    project: Project = Depends(updatable_project),
    projects: IProjectsService = Depends(projects_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_project,action=update:project",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> ProjectResponse:
    updated_project = await projects.update_project(project.id, updates)
    audit_log = props.new_event(
        f"project_id={updated_project.id},ref_id={user.id}",
        jsonable_encoder(updated_project),
    )
    tasks.add_task(audits.add_event, audit_log)
    return await projects.to_response(updated_project)


@router.delete(
    "/{project_id}",
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project not found"),
    ),
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["projects"],
)
async def delete_project(
    tasks: BackgroundTasks,
    user: User = Depends(current_user),
    project: Project = Depends(owned_project),
    projects: IProjectsService = Depends(projects_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_project,action=delete:project",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    await projects.delete_project(project.id)
    audit_log = props.new_event(f"project_id={project.id},ref_id={user.id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.get(
    "/{project_id}/members",
    response_model=ProjectMembersResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project not found"),
    ),
    status_code=status.HTTP_200_OK,
    tags=["projects"],
)
async def get_project_members(
    project: Project = Depends(viewable_project),
    projects: IProjectsService = Depends(projects_service),
    pages: Pagination = Depends(pagination),
) -> ProjectMembersResponse:
    members = await projects.get_members(project.id, pages.offset, pages.limit)
    return ProjectMembersResponse(
        items=[ProjectMemberResponse(**member.dict()) for member in members],
    )


@router.post(
    "/{project_id}/members",
    response_model=ProjectMemberResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project not found"),
    ),
    status_code=status.HTTP_200_OK,
    tags=["projects"],
)
async def add_project_member(
    tasks: BackgroundTasks,
    new_permission: NewPermissionRequest,
    user: User = Depends(current_user),
    project: Project = Depends(invitable_project),
    permissions: IPermissionsService = Depends(permissions_service),
    users: IUsersService = Depends(users_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=add_project_member,action=create:permission",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> ProjectMemberResponse:
    member_user = await users.get_by_id(new_permission.user_id)
    if not member_user:
        raise NotFoundError(message="User not found")

    if str(project.id) != str(new_permission.project_id):
        raise ConflictError(message="Project ids in path and payload differ")

    permission = await permissions.create_permission(
        NewPermissionRequest(**new_permission.dict()),
    )
    audit_log = props.new_event(
        f"permission_id={permission.id},ref_id={user.id}",
        jsonable_encoder(permission),
    )
    tasks.add_task(audits.add_event, audit_log)
    return ProjectMemberResponse(
        **dict(
            **member_user.dict(),
            permissions=PermissionResponse(**permission.dict()),
        ),
    )


@router.put(
    "/{project_id}/members/{user_id}",
    response_model=ProjectMemberResponse,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project or member not found"),
    ),
    status_code=status.HTTP_200_OK,
    tags=["projects"],
)
async def update_project_member(
    tasks: BackgroundTasks,
    user_id: Identifier,
    updated_permission: UpdatePermissionRequest,
    user: User = Depends(authenticated_user),
    project: Project = Depends(invitable_project),
    permissions: IPermissionsService = Depends(permissions_service),
    users: IUsersService = Depends(users_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_project_member,action=update:permission",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> ProjectMemberResponse:
    member_user = await users.get_by_id(user_id)
    old_permission = await permissions.find(project.id, user_id)
    permission = await permissions.update_permission(old_permission.id, updated_permission)
    audit_log = props.new_event(
        f"permission_id={permission.id},ref_id={user.id}",
        jsonable_encoder(permission),
    )
    tasks.add_task(audits.add_event, audit_log)
    return ProjectMemberResponse(
        **dict(
            **member_user.dict(),
            permissions=PermissionResponse(**permission.dict()),
        ),
    )


@router.delete(
    "/{project_id}/members/{user_id}",
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        (status.HTTP_404_NOT_FOUND, "Project or member not found"),
    ),
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["projects"],
)
async def delete_project_member(
    tasks: BackgroundTasks,
    user_id: Identifier,
    user: User = Depends(authenticated_user),
    project: Project = Depends(deletable_project),
    projects: IProjectsService = Depends(projects_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_project_member,action=delete:permission",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    deleted_permission = await projects.delete_member(project.id, user_id)
    audit_log = props.new_event(
        f"permission_id={deleted_permission.id},ref_id={user.id}",
        jsonable_encoder(deleted_permission),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
