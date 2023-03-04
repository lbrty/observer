from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.permissions import assert_viewable, assert_writable
from observer.common.types import Identifier, Role, SupportRecordSubject
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    people_service,
    permissions_service,
    pets_service,
    support_records_service,
)
from observer.entities.users import User
from observer.schemas.support_records import (
    NewSupportRecordRequest,
    SupportRecordResponse,
    UpdateSupportRecordRequest,
)
from observer.services.audit_logs import IAuditService
from observer.services.people import IPeopleService
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
    user: User = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    support_records: ISupportRecordsService = Depends(support_records_service),
    pets: IPetsService = Depends(pets_service),
    people: IPeopleService = Depends(people_service),
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
        subject_key = "person_id"
        await people.get_person(new_record.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(new_record.owner_id)

    support_record = await support_records.create_record(new_record)
    audit_log = props.new_event(
        f"{subject_key}={new_record.owner_id},ref_id={user.id}",
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
    user: User = Depends(
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
    user: User = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    people: IPeopleService = Depends(people_service),
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
        subject_key = "person_id"
        await people.get_person(updates.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(updates.owner_id)

    audit_log = props.new_event(
        f"{subject_key}={updates.owner_id},ref_id={user.id}",
        jsonable_encoder(support_record),
    )
    tasks.add_task(audits.add_event, audit_log)
    return SupportRecordResponse(**updated_support_record.dict())


@router.delete(
    "/{record_id}",
    status_code=status.HTTP_204_NO_CONTENT,
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
    user: User = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    people: IPeopleService = Depends(people_service),
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
) -> Response:
    support_record = await support_records.get_record(record_id)
    permission = await permissions.find(support_record.project_id, user.id)
    assert_writable(user, permission)
    await support_records.delete_record(record_id)

    if support_record.record_for == SupportRecordSubject.person:
        subject_key = "person_id"
        await people.get_person(support_record.owner_id)
    else:
        subject_key = "pet_id"
        await pets.get_pet(support_record.owner_id)

    audit_log = props.new_event(f"{subject_key}={support_record.owner_id},ref_id={user.id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
