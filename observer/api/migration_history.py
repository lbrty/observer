from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.permissions import (
    assert_can_see_private_info,
    assert_deletable,
    assert_viewable,
    assert_writable,
)
from observer.common.types import Identifier
from observer.components.audit import Props, Tracked
from observer.components.auth import authenticated_user
from observer.components.services import (
    audit_service,
    migrations_service,
    people_service,
    permissions_service,
)
from observer.entities.base import SomeUser
from observer.schemas.migration_history import (
    MigrationHistoryResponse,
    NewMigrationHistoryRequest,
)
from observer.services.audit_logs import IAuditService
from observer.services.migration_history import IMigrationService
from observer.services.people import IPeopleService
from observer.services.permissions import IPermissionsService

router = APIRouter(prefix="/migrations")


@router.post(
    "",
    response_model=MigrationHistoryResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    tags=["people", "migration history"],
)
async def create_migration_record(
    tasks: BackgroundTasks,
    new_record: NewMigrationHistoryRequest,
    user: SomeUser = Depends(authenticated_user),
    migrations: IMigrationService = Depends(migrations_service),
    audits: IAuditService = Depends(audit_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_migration_record,action=create:migration_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> MigrationHistoryResponse:
    """Create migration record"""
    permission = await permissions.find(new_record.project_id, user.id)
    assert_writable(user, permission)
    new_record = await migrations.add_record(new_record)
    audit_log = props.new_event(
        f"record_id={new_record.id},project_id={new_record.project_id},ref_id={user.ref_id}",
        jsonable_encoder(new_record),
    )
    tasks.add_task(audits.add_event, audit_log)
    return MigrationHistoryResponse(**new_record.dict())


@router.get(
    "/{record_id}",
    response_model=MigrationHistoryResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["people", "migration history"],
)
async def get_migration_record(
    record_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    people: IPeopleService = Depends(people_service),
    migrations: IMigrationService = Depends(migrations_service),
    permissions: IPermissionsService = Depends(permissions_service),
) -> MigrationHistoryResponse:
    """Get migration record"""
    migration_record = await migrations.get_record(record_id)
    idp_record = await people.get_person(migration_record.person_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    return MigrationHistoryResponse(**migration_record.dict())


@router.delete(
    "/{record_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["people", "migration history"],
)
async def delete_migration_record(
    tasks: BackgroundTasks,
    record_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    migrations: IMigrationService = Depends(migrations_service),
    audits: IAuditService = Depends(audit_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_migration_record,action=delete:migration_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    """Delete migration record"""
    migration_record = await migrations.get_record(record_id)
    permission = await permissions.find(migration_record.project_id, user.id)
    assert_deletable(user, permission)
    assert_can_see_private_info(user, permission)
    new_record = await migrations.delete_record(record_id)
    audit_log = props.new_event(
        f"record_id={new_record.id},project_id={new_record.project_id},ref_id={user.ref_id}",
        None,
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
