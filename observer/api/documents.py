from urllib.parse import quote

from fastapi import APIRouter, Depends, Response
from fastapi.responses import StreamingResponse
from starlette import status

from observer.common.permissions import assert_docs_readable, assert_viewable
from observer.common.types import Identifier, Role
from observer.components.auth import RequiresRoles
from observer.components.services import (
    documents_download,
    documents_service,
    permissions_service,
)
from observer.entities.base import SomeUser
from observer.schemas.documents import DocumentResponse
from observer.services.documents import IDocumentsService
from observer.services.downloads import DownloadHandler
from observer.services.permissions import IPermissionsService

router = APIRouter(prefix="/documents")


@router.get("/{doc_id}", tags=["documents"])
async def get_document() -> DocumentResponse:
    pass


@router.get(
    "/stream/{doc_id}",
    status_code=status.HTTP_200_OK,
    tags=["documents"],
)
async def stream_document(
    doc_id: Identifier,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.consultant, Role.staff]),
    ),
    permissions: IPermissionsService = Depends(permissions_service),
    documents: IDocumentsService = Depends(documents_service),
    downloads: DownloadHandler = Depends(documents_download),
) -> StreamingResponse:
    document = await documents.get_document(doc_id)
    permission = await permissions.find(document.project_id, user.id)
    assert_viewable(user, permission)
    assert_docs_readable(user, permission)
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


@router.delete("/{doc_id}", tags=["documents"])
async def delete_document() -> Response:
    pass
