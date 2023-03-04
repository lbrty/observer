from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked, client_ip
from observer.components.auth import RequiresRoles
from observer.components.pagination import pagination
from observer.components.services import (
    audit_service,
    auth_service,
    mailer,
    permissions_service,
    projects_service,
    users_service,
)
from observer.entities.users import User
from observer.schemas.pagination import Pagination
from observer.schemas.permissions import NewPermissionRequest
from observer.schemas.users import (
    NewUserRequest,
    UserInviteRequest,
    UserInviteResponse,
    UserInvitesResponse,
)
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.mailer import EmailMessage, IMailer
from observer.services.permissions import IPermissionsService
from observer.services.projects import IProjectsService
from observer.services.users import IUsersService
from observer.settings import settings

router = APIRouter(prefix="/invites")


@router.post(
    "",
    response_model=UserInviteResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["admin", "invites"],
)
async def create_invite(
    tasks: BackgroundTasks,
    invite_request: UserInviteRequest,
    user: User = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    users: IUsersService = Depends(users_service),
    auth: IAuthService = Depends(auth_service),
    audits: IAuditService = Depends(audit_service),
    permissions: IPermissionsService = Depends(permissions_service),
    projects: IProjectsService = Depends(projects_service),
    ip_address: str = Depends(client_ip),
    mail: IMailer = Depends(mailer),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_invite,action=create:invite",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> UserInviteResponse:
    random_password = await auth.gen_password()
    new_user = await users.create_user(
        NewUserRequest(
            email=invite_request.email,
            role=invite_request.role,
            password=random_password,
        ),
        is_active=False,
    )
    invite = await users.create_invite(new_user.id)
    link = f"https://{settings.app_domain}{settings.invite_url.format(code=invite.code)}"
    project_names = []
    if invite_request.permissions:
        for permission in invite_request.permissions:
            project = await projects.get_by_id(permission.project_id)
            project_names.append(project.name)
            permission_params = permission.dict()
            permission_params["user_id"] = new_user.id
            await permissions.create_permission(NewPermissionRequest(**permission_params))

    invited_projects = ""
    if project_names:
        invited_projects = "\nYou can access the following projects: " + ", ".join(project_names)

    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=new_user.email,
            from_email=settings.from_email,
            subject=settings.mfa_reset_subject,
            body=(
                f"You are invited to join {settings.invite_subject} please use the following {link}.{invited_projects}"
            ),
        ),
    )
    audit_log = props.new_event(
        f"action=send:invite,ref_id={user.id}",
        data=dict(new_user_id=new_user.id),
    )
    tasks.add_task(audits.add_event, audit_log)
    audit_log = props.new_event(
        f"ref_id={user.id}",
        data=dict(
            new_user_id=new_user.id,
            role=new_user.role.value,
            ip_address=ip_address,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return UserInviteResponse(**invite.dict())


@router.get(
    "",
    response_model=UserInvitesResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.staff]),
        )
    ],
    tags=["admin", "invites"],
)
async def get_invites(
    users: IUsersService = Depends(users_service),
    pages: Pagination = Depends(pagination),
) -> UserInvitesResponse:
    count, invites = await users.get_invites(pages)
    return UserInvitesResponse(total=count, items=invites)


@router.delete(
    "/{code}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["admin", "invites"],
)
async def delete_invite(
    code: Identifier,
    tasks: BackgroundTasks,
    delete_user: bool = Query(False, description="If set then related user will be deleted"),
    user: User = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    users: IUsersService = Depends(users_service),
    audits: IAuditService = Depends(audit_service),
    ip_address: str = Depends(client_ip),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_invite",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    invite = await users.get_invite(str(code), validate=False)
    await users.delete_invite(str(code))
    if delete_user:
        deleted_user = await users.delete_user(invite.user_id)
        audit_log = props.new_event(
            f"action=delete:user,ref_id={user.id}",
            data=dict(
                invite=jsonable_encoder(deleted_user),
                ip_address=ip_address,
            ),
        )
        tasks.add_task(audits.add_event, audit_log)

    audit_log = props.new_event(
        f"action=delete:invite,ref_id={user.id}",
        data=dict(
            invite=jsonable_encoder(invite),
            ip_address=ip_address,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
