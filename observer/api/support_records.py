from fastapi import APIRouter, BackgroundTasks, Depends
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.permissions import assert_viewable, assert_writable
from observer.common.types import Identifier, Role, SupportRecordSubject
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
    UpdateSupportRecordRequest,
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
        subject_key = "idp_id"
        await idp.get_idp(new_record.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(new_record.owner_id)

    support_record = await support_records.create_record(new_record)
    audit_log = props.new_event(
        f"{subject_key}={new_record.owner_id},ref_id={user.ref_id}",
        jsonable_encoder(support_record),
    )
    tasks.add_task(audits.add_event, audit_log)
    return SupportRecordResponse(**support_record.dict())


@router.get(
    "/{record_id}",
    response_model=SupportRecordResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["support records"],
)
async def get_support_record(
    record_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    support_records: ISupportRecordsService = Depends(support_records_service),
    permissions: IPermissionsService = Depends(permissions_service),
) -> SupportRecordResponse:
    support_record = await support_records.get_record(record_id)
    permission = await permissions.find(support_record.project_id, user.id)
    assert_viewable(user, permission)
    return SupportRecordResponse(**support_record.dict())


@router.put(
    "/{record_id}",
    response_model=SupportRecordResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["support records"],
)
async def update_support_record(
    tasks: BackgroundTasks,
    record_id: Identifier,
    updates: UpdateSupportRecordRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    idp: IIDPService = Depends(idp_service),
    pets: IPetsService = Depends(pets_service),
    support_records: ISupportRecordsService = Depends(support_records_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_support_record,action=update:support_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> SupportRecordResponse:
    support_record = await support_records.get_record(record_id)
    permission = await permissions.find(support_record.project_id, user.id)
    assert_writable(user, permission)
    updated_support_record = await support_records.update_record(record_id, updates)

    if updates.record_for == SupportRecordSubject.person:
        subject_key = "idp_id"
        await idp.get_idp(updates.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(updates.owner_id)

    audit_log = props.new_event(
        f"{subject_key}={updates.owner_id},ref_id={user.ref_id}",
        jsonable_encoder(support_record),
    )
    tasks.add_task(audits.add_event, audit_log)
    return SupportRecordResponse(**updated_support_record.dict())


@router.delete(
    "/{record_id}",
    response_model=SupportRecordResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["support records"],
)
async def delete_support_record(
    tasks: BackgroundTasks,
    record_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    idp: IIDPService = Depends(idp_service),
    pets: IPetsService = Depends(pets_service),
    support_records: ISupportRecordsService = Depends(support_records_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_support_record,action=delete:support_record",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> SupportRecordResponse:
    support_record = await support_records.get_record(record_id)
    permission = await permissions.find(support_record.project_id, user.id)
    assert_writable(user, permission)
    updated_support_record = await support_records.delete_record(record_id)

    if support_record.record_for == SupportRecordSubject.person:
        subject_key = "idp_id"
        await idp.get_idp(support_record.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(support_record.owner_id)

    audit_log = props.new_event(f"{subject_key}={support_record.owner_id},ref_id={user.ref_id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return SupportRecordResponse(**updated_support_record.dict())