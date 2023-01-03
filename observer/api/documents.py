from fastapi import APIRouter, UploadFile

from observer.schemas.documents import DocumentResponse

router = APIRouter(prefix="/documents")


@router.post("/upload", tags=["documents"])
async def upload_document(file: UploadFile) -> DocumentResponse:
    pass
