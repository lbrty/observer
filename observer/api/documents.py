from urllib.parse import quote

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from fastapi.responses import StreamingResponse
from starlette import status

from observer.api.exceptions import NotFoundError
from observer.common.exceptions import get_api_errors
from observer.common.permissions import (
    assert_deletable,
    assert_docs_readable,
    assert_viewable,
)
from observer.common.types import Role
from observer.components.audit import Props, Tracked
from observer.components.auth import RequiresRoles
from observer.components.documents import DocumentWithTests
from observer.components.services import (
    audit_service,
    documents_download,
    people_service,
    pets_service,
    storage_service,
)
from observer.entities.documents import Document
from observer.entities.users import User
from observer.schemas.documents import DocumentResponse
from observer.services.audit_logs import IAuditService
from observer.services.downloads import DownloadHandler
from observer.services.people import IPeopleService
from observer.services.pets import IPetsService
from observer.services.storage import IStorage

router = APIRouter(prefix="/documents")


@router.get(
    "/{doc_id}",
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["documents"],
)
async def get_document(
    document: Document = Depends(
        DocumentWithTests(
            assert_viewable,
            assert_docs_readable,
        )
    ),
) -> DocumentResponse:
    return DocumentResponse(**document.dict())


@router.get(
    "/stream/{doc_id}",
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["documents"],
)
async def stream_document(
    document: Document = Depends(
        DocumentWithTests(
            assert_viewable,
            assert_docs_readable,
        )
    ),
    downloads: DownloadHandler = Depends(documents_download),
) -> StreamingResponse:
    content_disposition_filename = quote(document.name)
    if content_disposition_filename != document.name:
        content_disposition = "{}; filename*=utf-8''{}".format(document.mimetype, content_disposition_filename)
    else:
        content_disposition = '{}; filename="{}"'.format(document.mimetype, document.name)

    return StreamingResponse(
        downloads.stream(document),
        media_type=document.mimetype,
        headers={"content-disposition": content_disposition},
    )


@router.delete(
    "/{doc_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    dependencies=[
        Depends(
            RequiresRoles([Role.admin, Role.consultant, Role.staff]),
        )
    ],
    tags=["documents"],
)
async def delete_document(
    tasks: BackgroundTasks,
    user: User = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    document: Document = Depends(
        DocumentWithTests(
            assert_viewable,
            assert_deletable,
            assert_docs_readable,
        )
    ),
    pets: IPetsService = Depends(pets_service),
    people: IPeopleService = Depends(people_service),
    storage: IStorage = Depends(storage_service),
    audits: IAuditService = Depends(audit_service),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_document,action=delete:document",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    subject_key = None
    try:
        await pets.get_pet(document.owner_id)
        subject_key = "pet_id"
    except NotFoundError:
        pass

    if not subject_key:
        try:
            await people.get_person(document.owner_id)
            subject_key = "person_id"
        except NotFoundError:
            pass

    tasks.add_task(storage.delete_path, document.path)
    audit_log = props.new_event(f"{subject_key}={document.owner_id},ref_id={user.id}", None)
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
