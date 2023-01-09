from fastapi import APIRouter, BackgroundTasks, Depends
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.permissions import assert_writable
from observer.common.types import Role, SupportRecordSubject
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    idp_service,
    permissions_service,
    pets_service,
    support_records_service,
)
from observer.entities.base import SomeUser
from observer.schemas.support_records import (
    NewSupportRecordRequest,
    SupportRecordResponse,
)
from observer.services.audit_logs import IAuditService
from observer.services.idp import IIDPService
from observer.services.permissions import IPermissionsService
from observer.services.pets import IPetsService
from observer.services.support_records import ISupportRecordsService

router = APIRouter(prefix="/support-records")


@router.post(
    "",
    response_model=SupportRecordResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["support records"],
)
async def create_support_record(
    tasks: BackgroundTasks,
    new_record: NewSupportRecordRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    support_records: ISupportRecordsService = Depends(support_records_service),
    pets: IPetsService = Depends(pets_service),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_support_record,action=create:support_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> SupportRecordResponse:
    permission = await permissions.find(new_record.project_id, user.id)
    assert_writable(user, permission)

    if new_record.record_for == SupportRecordSubject.person:
        await idp.get_idp(new_record.owner_id)

    if new_record.record_for == SupportRecordSubject.pet:
        await pets.get_pet(new_record.owner_id)

    support_record = await support_records.create_record(new_record)
    return SupportRecordResponse(**support_record.dict())
