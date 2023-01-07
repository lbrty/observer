from fastapi import APIRouter, Response, UploadFile
from fastapi.responses import StreamingResponse

from observer.schemas.documents import DocumentResponse

router = APIRouter(prefix="/documents")


@router.get("/{doc_id}", tags=["documents"])
async def get_document(file: UploadFile) -> DocumentResponse:
    pass


@router.get("/{doc_id}/stream", tags=["documents"])
async def stream_document(file: UploadFile) -> StreamingResponse:
    pass


@router.delete("/{doc_id}", tags=["documents"])
async def delete_document(file: UploadFile) -> Response:
    pass
