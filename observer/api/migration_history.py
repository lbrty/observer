from fastapi import APIRouter, BackgroundTasks, Depends
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.permissions import assert_writable
from observer.components.audit import Props, Tracked
from observer.components.auth import authenticated_user
from observer.components.services import (
    audit_service,
    idp_service,
    migrations_service,
    permissions_service,
)
from observer.entities.base import SomeUser
from observer.schemas.migration_history import (
    MigrationHistoryResponse,
    NewMigrationHistoryRequest,
)
from observer.services.audit_logs import IAuditService
from observer.services.idp import IIDPService
from observer.services.migration_history import IMigrationService
from observer.services.permissions import IPermissionsService

router = APIRouter(prefix="/migrations")


@router.post(
    "",
    response_model=MigrationHistoryResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["idp", "migration", "history"],
)
async def create_migration_record(
    tasks: BackgroundTasks,
    new_record: NewMigrationHistoryRequest,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
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
    idp_record = await idp.get_idp(new_record.idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_writable(user, permission)
    new_record = await migrations.add_record(new_record)
    audit_log = props.new_event(
        f"record_id={new_record.id},ref_id={user.ref_id}",
        jsonable_encoder(new_record),
    )
    tasks.add_task(audits.add_event, audit_log)
    return MigrationHistoryResponse(**new_record.dict())
