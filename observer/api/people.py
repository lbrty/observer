import os
from typing import List

from fastapi import APIRouter, BackgroundTasks, Depends, Response, UploadFile
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.api.exceptions import ConflictError
from observer.common.exceptions import get_api_errors
from observer.common.permissions import (
    assert_can_see_private_info,
    assert_deletable,
    assert_docs_readable,
    assert_updatable,
    assert_viewable,
    assert_writable,
)
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked, client_ip
from observer.components.auth import RequiresRoles, authenticated_user
from observer.components.services import (
    audit_service,
    documents_service,
    documents_upload,
    family_service,
    idp_service,
    migrations_service,
    permissions_service,
    secrets_service,
    storage_service,
    world_service,
)
from observer.entities.base import SomeUser
from observer.entities.people import PersonalInfo
from observer.schemas.documents import DocumentResponse, NewDocumentRequest
from observer.schemas.family_members import (
    FamilyMemberResponse,
    NewFamilyMemberRequest,
    UpdateFamilyMemberRequest,
)
from observer.schemas.people import (
    IDPResponse,
    NewIDPRequest,
    PersonalInfoResponse,
    UpdateIDPRequest,
)
from observer.schemas.migration_history import FullMigrationHistoryResponse
from observer.schemas.world import PlaceResponse
from observer.services.audit_logs import IAuditService
from observer.services.documents import IDocumentsService
from observer.services.family_members import IFamilyService
from observer.services.people import IIDPService
from observer.services.migration_history import IMigrationService
from observer.services.permissions import IPermissionsService
from observer.services.secrets import ISecretsService
from observer.services.storage import IStorage
from observer.services.uploads import UploadHandler
from observer.services.world import IWorldService
from observer.settings import settings

router = APIRouter(prefix="/people")


@router.post(
    "",
    response_model=IDPResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
    ),
    tags=["idp", "people"],
)
async def create_idp(
    tasks: BackgroundTasks,
    new_idp: NewIDPRequest,
    user: SomeUser = Depends(authenticated_user),
    audits: IAuditService = Depends(audit_service),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_idp,action=create:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> IDPResponse:
    permission = await permissions.find(new_idp.project_id, user.id)
    assert_writable(user, permission)
    person = await idp.create_idp(new_idp)
    audit_log = props.new_event(f"person_id={person.id},ref_id={user.ref_id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return IDPResponse(**person.dict())


@router.get(
    "/{idp_id}",
    response_model=IDPResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_404_NOT_FOUND,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people"],
)
async def get_idp(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    secrets: ISecretsService = Depends(secrets_service),
) -> IDPResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    idp_record = await secrets.anonymize_idp(idp_record)
    return IDPResponse(**idp_record.dict())


@router.get(
    "/{idp_id}/personal-info",
    response_model=PersonalInfoResponse,
    response_model_exclude_none=True,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people"],
)
async def get_personal_info(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    secrets: ISecretsService = Depends(secrets_service),
) -> PersonalInfoResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_can_see_private_info(user, permission)
    pi = await secrets.decrypt_personal_info(
        PersonalInfo(
            email=idp_record.email,
            phone_number=idp_record.phone_number,
            phone_number_additional=idp_record.phone_number_additional,
        ),
    )
    pi.full_name = idp_record.full_name
    pi.sex = idp_record.sex
    pi.pronoun = idp_record.pronoun
    return PersonalInfoResponse(**pi.dict())


@router.get(
    "/{idp_id}/migration-records",
    response_model=List[FullMigrationHistoryResponse],
    response_model_exclude_none=True,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people", "migration", "history"],
)
async def get_person_migration_records(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    migrations: IMigrationService = Depends(migrations_service),
    world: IWorldService = Depends(world_service),
) -> List[FullMigrationHistoryResponse]:
    """Get migration records for a person"""
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_can_see_private_info(user, permission)
    records = await migrations.get_idp_records(idp_id)
    result = []
    for record in records:
        migration_record = FullMigrationHistoryResponse(**record.dict())
        if record.from_place_id:
            place = await world.get_place(record.from_place_id)
            migration_record.from_place = PlaceResponse(**place.dict())

        if record.current_place_id:
            place = await world.get_place(record.current_place_id)
            migration_record.current_place = PlaceResponse(**place.dict())
        result.append(migration_record)

    return result


@router.put(
    "/{idp_id}",
    response_model=IDPResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people"],
)
async def update_idp(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    idp_updates: UpdateIDPRequest,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_idp,action=update:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> IDPResponse:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_updatable(user, permission)
    updated_idp = await idp.update_idp(idp_id, idp_updates)
    audit_log = props.new_event(
        f"person_id={updated_idp.id},ref_id={user.ref_id}",
        jsonable_encoder(
            updated_idp.dict(
                exclude_none=True,
                exclude={"email", "phone_number", "phone_number_additional"},
            )
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return IDPResponse(**updated_idp.dict())


@router.delete(
    "/{idp_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people"],
)
async def delete_idp(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
    storage: IStorage = Depends(storage_service),
    audits: IAuditService = Depends(audit_service),
    ip_address: str = Depends(client_ip),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_idp,action=delete:idp",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_deletable(user, permission)
    assert_docs_readable(user, permission)
    deleted_idp = await idp.delete_idp(idp_id)
    idp_documents = await documents.get_by_owner_id(idp_id)
    document_ids = [str(doc.id) for doc in idp_documents]
    await documents.bulk_delete(document_ids)
    full_path = os.path.join(settings.documents_path, str(idp_id))
    audit_log = props.new_event(f"person_id={deleted_idp.id},ref_id={user.ref_id}", dict(ip_address=ip_address))
    tasks.add_task(storage.delete_path, full_path)
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/{idp_id}/document",
    response_model=DocumentResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "documents"],
)
async def idp_upload_document(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    file: UploadFile,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
    uploads: UploadHandler = Depends(documents_upload),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=idp_upload_document,action=create:document",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> DocumentResponse:
    pet = await idp.get_idp(idp_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_deletable(user, permission)
    assert_docs_readable(user, permission)
    save_to = os.path.join(settings.documents_path, str(idp_id))
    size, sealed_file = await uploads.process_upload(file, save_to)
    document = await documents.create_document(
        sealed_file.encryption_key,
        NewDocumentRequest(
            name=file.filename,
            size=size,
            path=sealed_file.path,
            mimetype=file.content_type,
            owner_id=idp_id,
            project_id=pet.project_id,
        ),
    )
    audit_log = props.new_event(
        f"idp_id={pet.id},ref_id={user.ref_id}",
        jsonable_encoder(
            document,
            exclude={"id", "encryption_key"},
            exclude_none=True,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return DocumentResponse(**document.dict())


@router.get(
    "/{idp_id}/documents",
    response_model=List[DocumentResponse],
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "documents"],
)
async def idp_get_documents(
    idp_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    idp: IIDPService = Depends(idp_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
) -> List[DocumentResponse]:
    pet = await idp.get_idp(idp_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_viewable(user, permission)
    assert_docs_readable(user, permission)
    docs = await documents.get_by_owner_id(idp_id)
    return [DocumentResponse(**doc.dict()) for doc in docs]


@router.post(
    "/{idp_id}/family-members",
    response_model=FamilyMemberResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["idp", "people", "family members"],
)
async def add_persons_family_member(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    new_member: NewFamilyMemberRequest,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    family: IFamilyService = Depends(family_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=add_persons_family_member,action=create:family_member",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> FamilyMemberResponse:
    """Add family member for a person"""
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    if idp_record.project_id != new_member.project_id:
        raise ConflictError(message="Project ID is not the same as in request body")

    assert_viewable(user, permission)
    assert_can_see_private_info(user, permission)
    member = await family.add_member(new_member)
    audit_log = props.new_event(
        f"idp_id={idp_id},project_id={new_member.project_id},member_id={member.id},ref_id={user.ref_id}",
        jsonable_encoder(
            member,
            exclude={"id"},
            exclude_none=True,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return FamilyMemberResponse(**member.dict())


@router.get(
    "/{idp_id}/family-members",
    response_model=List[FamilyMemberResponse],
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people", "family members"],
)
async def get_persons_family_members(
    idp_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    family: IFamilyService = Depends(family_service),
    permissions: IPermissionsService = Depends(permissions_service),
) -> List[FamilyMemberResponse]:
    """Get family members for a person"""
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    assert_can_see_private_info(user, permission)
    members = await family.get_by_person(idp_id)
    return [FamilyMemberResponse(**member.dict()) for member in members]


@router.put(
    "/{idp_id}/family-members/{member_id}",
    response_model=FamilyMemberResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people", "family members"],
)
async def update_persons_family_member(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    member_id: Identifier,
    updates: UpdateFamilyMemberRequest,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    family: IFamilyService = Depends(family_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_persons_family_member,action=update:family_member",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> FamilyMemberResponse:
    """Add family member for a person"""
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    assert_can_see_private_info(user, permission)
    member = await family.update_member(member_id, updates)
    audit_log = props.new_event(
        f"idp_id={idp_id},project_id={member.project_id},member_id={member.id},ref_id={user.ref_id}",
        jsonable_encoder(
            member,
            exclude={"id"},
            exclude_none=True,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return FamilyMemberResponse(**member.dict())


@router.delete(
    "/{idp_id}/family-members/{member_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["idp", "people", "family members"],
)
async def delete_persons_family_member(
    tasks: BackgroundTasks,
    idp_id: Identifier,
    member_id: Identifier,
    user: SomeUser = Depends(authenticated_user),
    idp: IIDPService = Depends(idp_service),
    family: IFamilyService = Depends(family_service),
    permissions: IPermissionsService = Depends(permissions_service),
    audits: IAuditService = Depends(audit_service),
    ip_address: str = Depends(client_ip),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_persons_family_member,action=delete:family_member",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    """Delete family member for a person"""
    idp_record = await idp.get_idp(idp_id)
    permission = await permissions.find(idp_record.project_id, user.id)
    assert_viewable(user, permission)
    assert_can_see_private_info(user, permission)
    member = await family.delete_member(member_id)
    audit_log = props.new_event(
        f"idp_id={idp_id},project_id={member.project_id},member_id={member.id},ref_id={user.ref_id}",
        dict(ip_address=ip_address),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
