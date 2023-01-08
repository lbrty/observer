import os
from typing import List

from fastapi import APIRouter, Depends, Response, UploadFile
from fastapi.encoders import jsonable_encoder
from starlette import status
from starlette.background import BackgroundTasks

from observer.common.permissions import (
    assert_deletable,
    assert_docs_readable,
    assert_updatable,
    assert_viewable,
    assert_writable,
)
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    documents_service,
    documents_upload,
    permissions_service,
    pets_service,
    storage_service,
)
from observer.components.validate import ContentLengthLimit
from observer.entities.base import SomeUser
from observer.schemas.documents import DocumentResponse, NewDocumentRequest
from observer.schemas.pets import NewPetRequest, PetResponse, UpdatePetRequest
from observer.services.audit_logs import IAuditService
from observer.services.documents import IDocumentsService
from observer.services.permissions import IPermissionsService
from observer.services.pets import IPetsService
from observer.services.storage import IStorage
from observer.services.uploads import UploadHandler
from observer.settings import settings

router = APIRouter(prefix="/pets")


@router.post(
    "",
    response_model=PetResponse,
    status_code=status.HTTP_201_CREATED,
    tags=["pets"],
)
async def create_pet(
    tasks: BackgroundTasks,
    new_pet: NewPetRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    audits: IAuditService = Depends(audit_service),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_pet,action=create:pet",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> PetResponse:
    permission = await permissions.find(new_pet.project_id, user.id)
    assert_writable(user, permission)
    pet = await pets.create_pet(new_pet)
    audit_log = props.new_event(f"pet_id={pet.id},ref_id={user.ref_id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return PetResponse(**pet.dict())


@router.get(
    "/{pet_id}",
    response_model=PetResponse,
    status_code=status.HTTP_200_OK,
    tags=["pets"],
)
async def get_pet(
    pet_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
) -> PetResponse:
    pet = await pets.get_pet(pet_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_viewable(user, permission)
    return PetResponse(**pet.dict())


@router.put(
    "/{pet_id}",
    response_model=PetResponse,
    status_code=status.HTTP_200_OK,
    tags=["pets"],
)
async def update_pet(
    tasks: BackgroundTasks,
    pet_id: Identifier,
    updates: UpdatePetRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    audits: IAuditService = Depends(audit_service),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=update_pet,action=update:pet",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> PetResponse:
    pet = await pets.get_pet(pet_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_updatable(user, permission)
    updated_pet = await pets.update_pet(pet_id, updates)
    audit_log = props.new_event(
        f"pet_id={pet.id},ref_id={user.ref_id}",
        jsonable_encoder(updated_pet, exclude={"id"}, exclude_none=True),
    )
    tasks.add_task(audits.add_event, audit_log)
    return PetResponse(**updated_pet.dict())


@router.delete(
    "/{pet_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    tags=["pets"],
)
async def delete_pet(
    tasks: BackgroundTasks,
    pet_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    audits: IAuditService = Depends(audit_service),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
    storage: IStorage = Depends(storage_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_pet,action=delete:pet",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    pet = await pets.get_pet(pet_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_deletable(user, permission)
    assert_docs_readable(user, permission)
    pet_documents = await documents.get_by_owner_id(pet_id)
    document_ids = [str(doc.id) for doc in pet_documents]
    await documents.bulk_delete(document_ids)
    full_path = os.path.join(settings.documents_path, str(pet_id))
    deleted_pet = await pets.delete_pet(pet_id)
    audit_log = props.new_event(
        f"pet_id={pet.id},ref_id={user.ref_id}",
        jsonable_encoder(deleted_pet, exclude={"id"}, exclude_none=True),
    )
    tasks.add_task(audits.add_event, audit_log)
    await storage.delete_path(full_path)
    # tasks.add_task(storage.delete_path, full_path)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/{pet_id}/document",
    response_model=DocumentResponse,
    status_code=status.HTTP_201_CREATED,
    dependencies=[
        Depends(ContentLengthLimit(settings.max_upload_size)),
    ],
    tags=["pets", "documents"],
)
async def pet_upload_document(
    tasks: BackgroundTasks,
    pet_id: Identifier,
    file: UploadFile,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    audits: IAuditService = Depends(audit_service),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
    uploads: UploadHandler = Depends(documents_upload),
    props: Props = Depends(
        Tracked(
            tag="endpoint=pet_upload_document,action=create:document",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> DocumentResponse:
    pet = await pets.get_pet(pet_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_deletable(user, permission)
    assert_docs_readable(user, permission)
    save_to = os.path.join(settings.documents_path, str(pet_id))
    size, sealed_file = await uploads.process_upload(file, save_to)
    document = await documents.create_document(
        sealed_file.encryption_key,
        NewDocumentRequest(
            name=file.filename,
            size=size,
            path=sealed_file.path,
            mimetype=file.content_type,
            owner_id=pet_id,
            project_id=pet.project_id,
        ),
    )
    audit_log = props.new_event(
        f"pet_id={pet.id},ref_id={user.ref_id}",
        jsonable_encoder(
            document,
            exclude={"id", "encryption_key"},
            exclude_none=True,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return DocumentResponse(**document.dict())


@router.get(
    "/{pet_id}/documents",
    response_model=List[DocumentResponse],
    status_code=status.HTTP_200_OK,
    tags=["pets", "documents"],
)
async def pet_get_documents(
    pet_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    pets: IPetsService = Depends(pets_service),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
) -> List[DocumentResponse]:
    pet = await pets.get_pet(pet_id)
    permission = await permissions.find(pet.project_id, user.id)
    assert_docs_readable(user, permission)
    docs = await documents.get_by_owner_id(pet_id)
    return [DocumentResponse(**doc.dict()) for doc in docs]
