from fastapi import APIRouter, BackgroundTasks, Depends
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    permissions_service,
    pets_service,
)
from observer.entities.base import SomeUser
from observer.schemas.pets import PetResponse
from observer.schemas.support_records import NewSupportRecordRequest
from observer.services.audit_logs import IAuditService
from observer.services.permissions import IPermissionsService
from observer.services.pets import IPetsService

router = APIRouter(prefix="/support-records")


@router.post(
    "",
    response_model=PetResponse,
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
    audits: IAuditService = Depends(audit_service),
    support_records: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_support_record,action=create:support_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> PetResponse:
    ...
